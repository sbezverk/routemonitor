package radix

import (
	"fmt"
	"net"
	"strings"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bgp"
)

// Tree is the interface to manage tree related functions
type Tree interface {
	Add([]byte, int, string, *bgp.BaseAttributes)
	Check([]byte, int) bool
	GetAll() []string
}

type prefix struct {
	value  []byte
	length int
	attrs  map[string]*bgp.BaseAttributes
}

var _ Tree = &tree{}

type node struct {
	parent *node
	right  *node
	left   *node
	prefix *prefix
}

type msg struct {
	op     operation
	value  []byte
	length int
	peer   string
	attr   *bgp.BaseAttributes
}

type tree struct {
	root     *node
	treeCh   chan msg
	resultCh chan interface{}
}

type operation int

const (
	addOp operation = iota
	delOp
	checkOp
)

// treeManager is function sceduling operations on the tree, it is used to prevent any concurrency issues
func (t *tree) treeManager() {
	for {
		select {
		case msg := <-t.treeCh:
			switch msg.op {
			case addOp:
				t.add(msg.value, msg.length, msg.peer, msg.attr)
			case delOp:
			case checkOp:
				t.check(msg.value, msg.length)
			}
		default:
		}
	}
}

// Add is externally available methor to add a route into the tree
func (t *tree) Add(b []byte, l int, peer string, attr *bgp.BaseAttributes) {
	t.treeCh <- msg{
		op:     addOp,
		value:  b,
		length: l,
		peer:   peer,
		attr:   attr,
	}
}

// Check verifies if specified prefix stored in the tree
func (t *tree) Check(b []byte, l int) bool {
	t.treeCh <- msg{
		op:     checkOp,
		value:  b,
		length: l,
	}
	r := <-t.resultCh
	_, ok := r.(bool)
	if !ok {
		return false
	}

	return r.(bool)
}

func (t *tree) check(b []byte, l int) {
	//	if n.root == nil {
	//		n.root = &node{}
	//	}
	glog.Infof("Checking for prefix: %+v length: %d", b, l)
	cnode := t.root
	v := NewNodeValue()
	v.LoadNodeValue(b)
	i := 0
	for d := range v.BitRanger() {
		if d {
			if cnode.right == nil {
				t.resultCh <- false
				return
			}
			cnode = cnode.right
		} else {
			if cnode.left == nil {
				t.resultCh <- false
				return
			}
			cnode = cnode.left
		}
		if i >= l {
			break
		}
		i++
	}
	if cnode.prefix == nil {
		t.resultCh <- false
		return
	}
	t.resultCh <- true
}

func (t *tree) add(b []byte, l int, peer string, attr *bgp.BaseAttributes) {
	//	if n.root == nil {
	//		n.root = &node{}
	//	}
	cnode := t.root
	v := NewNodeValue()
	v.LoadNodeValue(b)
	i := 0
	for t := range v.BitRanger() {
		if t {
			if cnode.right == nil {
				cnode.right = &node{}
				cnode.right.parent = cnode
			}
			cnode = cnode.right
		} else {
			if cnode.left == nil {
				cnode.left = &node{}
				cnode.left.parent = cnode
			}
			cnode = cnode.left
		}
		if i >= l {
			break
		}
		i++
	}
	if cnode.prefix == nil {
		cnode.prefix = &prefix{
			length: l,
			attrs:  make(map[string]*bgp.BaseAttributes),
		}
	}
	oattr, ok := cnode.prefix.attrs[peer]
	if !ok {
		cnode.prefix.attrs[peer] = attr
	} else {
		// Anpther advertisement from the same peer for this prefix
		// compare BaseAttributes hash
		if strings.Compare(oattr.BaseAttrHash, attr.BaseAttrHash) != 0 {
			// Changes in attributes detected, saving moe recent attributes
			cnode.prefix.attrs[peer] = attr
			// If change in attributes is tracking, here the signal should be raised to notify about the change
		}
	}
}

func setBit(b []byte, n int) error {
	if n >= len(b)*8 {
		return fmt.Errorf("invalid bit bit %d, slice is only %d bits long", n, len(b)*8)
	}
	i := n / 8
	b[i] += 0x80 >> (n % 8)

	return nil
}

func clearBit(b []byte, n int) {
	i := n / 8
	if n%8 == 0 && i != 0 {
		i--
	}
	mask := ^(1 << (8 - (n - (i * 8))))
	b[i] &= byte(mask)
}

func copyBits(d, s []byte, n int) {
	copy(d, s[:len(d)])
	i := n / 8
	if i != 0 {
		i--
	}
	for m := 0; m < (n % 8); m++ {
		mask := ^(1 << m)
		d[i] &= byte(mask)
	}
}

func (t *tree) processNode(cnode *node, onode *node, up bool, bit int, c chan *prefix, p []byte) (*node, bool, int) {
	if cnode == nil {
		close(c)
		return nil, up, bit
	}
	// Check for prefixes only on the way down of the tree
	if cnode.prefix != nil && !up {
		pr := &prefix{
			length: cnode.prefix.length,
		}
		pl := cnode.prefix.length / 8
		if cnode.prefix.length%8 != 0 {
			pl++
		}
		pr.value = make([]byte, pl)
		copyBits(pr.value, p, cnode.prefix.length)
		c <- pr
	}
	if up {
		// Direction is back up
		// On direction up, make sense to check only the right side, as left has already been traversed
		// on the way down.
		// If current node's right is not nil and not the same as the old node, meaning, the walk did not
		// come back from the right child, then it is a new/unvisited branch.
		if cnode.right != nil && cnode.right != onode {
			setBit(p, bit)
			up = false
			bit++
			return t.processNode(cnode.right, cnode, up, bit, c, p)
		}
		// No instantiated children nodes, going back up
		up = true
		clearBit(p, bit)
		bit--
		return t.processNode(cnode.parent, cnode, up, bit, c, p)
	}
	if cnode.left != nil {
		up = false
		bit++
		return t.processNode(cnode.left, cnode, up, bit, c, p)
	}
	if cnode.right != nil {
		setBit(p, bit)
		up = false
		bit++
		return t.processNode(cnode.right, cnode, up, bit, c, p)
	}
	// No instantiated children nodes, going back up
	up = true
	clearBit(p, bit)
	bit--
	return t.processNode(cnode.parent, cnode, up, bit, c, p)
}

func (t *tree) GetAll() []string {
	routes := make([]string, 0)
	c := make(chan *prefix)
	p := make([]byte, 4)
	go t.processNode(t.root, t.root.parent, false, 0, c, p)
	for p := range c {
		pr := make([]byte, 4)
		copy(pr, p.value)
		routes = append(routes, fmt.Sprintf("%s/%d", net.IP(pr).To4().String(), p.length))
	}

	return routes
}

// NewTree returns a new instance of the tree
func NewTree() Tree {
	t := &tree{
		treeCh:   make(chan msg),
		resultCh: make(chan interface{}),
		root: &node{
			parent: nil,
		},
	}
	// Starting Tree Manager
	go t.treeManager()

	return t
}

// NodeValue defines interface with methods to operate with Node Value
type NodeValue interface {
	LoadNodeValue(b []byte)
	BitRanger() chan bool
}

type nodeValue struct {
	value []byte
}

// NewNodeValue instantiates a new instance of node value used to build the tree
func NewNodeValue() NodeValue {
	return &nodeValue{}
}
func (nv *nodeValue) LoadNodeValue(b []byte) {
	nv.value = make([]byte, len(b))
	copy(nv.value, b)
}

func (nv *nodeValue) BitRanger() chan bool {
	c := make(chan bool)
	go func() {
		for i := 0; i < len(nv.value); i++ {
			s := 0x80
			for y := 7; y >= 0; y-- {
				r := nv.value[i]&uint8(s) == uint8(s)
				s >>= 1
				c <- r
			}
		}
		close(c)
		return
	}()
	return c
}
