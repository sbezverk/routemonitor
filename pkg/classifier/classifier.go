package classifier

import (
	"fmt"
	"net"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/base"
	"github.com/sbezverk/gobmp/pkg/bgp"
	"github.com/sbezverk/gobmp/pkg/bmp"
	"github.com/sbezverk/routemonitor/pkg/radix"
)

// RouteType defines type of routes which routemontor supports
type RouteType int

const (
	// UnicastIPv4 unicast ipv4 type of route
	UnicastIPv4 RouteType = iota
	// UnicastIPv6 unicast ipv6 type of route
	UnicastIPv6
	// VPNv4 l3 vpnv4 type of route
	VPNv4
	// VPNv6 l3 vpnv6 type of route
	VPNv6
)

// NLRI defines method to classify a type of NLRI found in BMP Route Monitor message
type NLRI interface {
	Classify(bmp.Message)
	Check(RouteType, []byte, int) bool
	Delete(RouteType, []byte, int, string) error
	GetAll(RouteType) []string
	Monitor(RouteType, string, []byte, int, chan struct{}) error
	Unmonitor(RouteType, string, []byte, int, chan struct{})
}

var _ NLRI = &nlri{}

type nlri struct {
	uipv4 radix.Tree
	vpnv4 radix.Tree
	uipv6 radix.Tree
	vpnv6 radix.Tree
}

func (n *nlri) Check(t RouteType, b []byte, l int) bool {
	switch t {
	case UnicastIPv4:
		return n.uipv4.Check(b, l)
	case UnicastIPv6:
		return n.uipv6.Check(b, l)
	case VPNv4:
		return n.vpnv4.Check(b, l)
	case VPNv6:
		return n.vpnv6.Check(b, l)
	}

	return false
}

func (n *nlri) Monitor(t RouteType, id string, b []byte, l int, c chan struct{}) error {
	switch t {
	case UnicastIPv4:
		return n.uipv4.Monitor(id, b, l, c)
	case UnicastIPv6:
		return n.uipv6.Monitor(id, b, l, c)
	case VPNv4:
		return n.vpnv4.Monitor(id, b, l, c)
	case VPNv6:
		return n.vpnv6.Monitor(id, b, l, c)
	}

	return fmt.Errorf("not supported table")
}

func (n *nlri) Unmonitor(t RouteType, id string, b []byte, l int, c chan struct{}) {
	switch t {
	case UnicastIPv4:
		n.uipv4.Unmonitor(id, b, l, c)
		return
	case UnicastIPv6:
		n.uipv6.Unmonitor(id, b, l, c)
		return
	case VPNv4:
		n.vpnv4.Unmonitor(id, b, l, c)
		return
	case VPNv6:
		n.vpnv6.Unmonitor(id, b, l, c)
		return
	}

	return
}

func (n *nlri) Delete(t RouteType, b []byte, l int, peer string) error {
	switch t {
	case UnicastIPv4:
		return n.uipv4.Delete(b, l, peer)
	case UnicastIPv6:
		return n.uipv6.Delete(b, l, peer)
	case VPNv4:
		return n.vpnv4.Delete(b, l, peer)
	case VPNv6:
		return n.vpnv6.Delete(b, l, peer)
	}

	return fmt.Errorf("non supported route type")
}

func (n *nlri) Classify(msg bmp.Message) {
	// If PeerHeader is nil, no point look further
	if msg.PeerHeader == nil {
		return
	}
	var peer string
	switch msg.PeerHeader.PeerType {
	case 0:
		peer += "Global peer: "
	case 1:
		peer += "VPN peer: "
	case 2:
		peer += "Local peer: "
	}
	peer += msg.PeerHeader.GetPeerDistinguisherString() + " "
	if msg.PeerHeader.FlagV {
		peer += fmt.Sprintf("%s", net.IP(msg.PeerHeader.PeerAddress).To16().String())
	} else {
		peer += fmt.Sprintf("%s", net.IP(msg.PeerHeader.PeerAddress[len(msg.PeerHeader.PeerAddress)-4:]).To4().String())
	}
	m, ok := msg.Payload.(*bmp.RouteMonitor)
	if !ok {
		return
	}
	u := m.Update
	glog.V(5).Infof("Processing BMP from peer: %s peer hash: %s", peer, msg.PeerHeader.GetPeerHash())
	if len(u.NLRI) != 0 {
		n.processBGPNLRI(msg.PeerHeader.GetPeerHash(), u.BaseAttributes, u.NLRI)
	}
	if len(u.WithdrawnRoutes) != 0 {
		n.processBGPWithdraw(msg.PeerHeader.GetPeerHash(), u.BaseAttributes, u.WithdrawnRoutes)
	}
	if mpreach, _ := u.GetMPReachNLRI(); mpreach != nil {
		n.processMPReach(msg.PeerHeader.GetPeerHash(), u.BaseAttributes, mpreach)
	}
	if mpunreach, _ := u.GetMPUnReachNLRI(); mpunreach != nil {
		n.processMPUnReach(msg.PeerHeader.GetPeerHash(), u.BaseAttributes, mpunreach)
	}
}

func (n *nlri) processBGPNLRI(peer string, attr *bgp.BaseAttributes, routes []base.Route) {
	glog.V(5).Info("Message with bgp rfc4271 nlri")
	for _, r := range routes {
		// glog.Infof("><SB> BGP NLRI Prefix: %+v length: %d", r.Prefix, r.Length)
		n.uipv4.Add(r.Prefix, int(r.Length), peer, attr)
	}
}

func (n *nlri) processBGPWithdraw(peer string, attr *bgp.BaseAttributes, routes []base.Route) {
	glog.V(5).Info("Message with bgp rfc4271 withdraw")
	// fmt.Printf("Update carries rfc4271 Withdraw routes\n")
	// fmt.Printf("Peer: %s Attributes: %+vRoutes: %+v", peer, attr, routes)
}
func (n *nlri) processMPReach(peer string, attr *bgp.BaseAttributes, mpreach bgp.MPNLRI) {
	glog.V(5).Infof("Message with bgp mp_reach nlri, afi/safi code: %d", mpreach.GetAFISAFIType())
	if unicast, err := mpreach.GetNLRIUnicast(); err == nil {
		glog.V(5).Infof("Message with bgp mp_reach nlri unicast")
		// fmt.Printf("Peer: %s Attributes: %+v Unicast Routes: %+v", peer, attr, unicast.NLRI)
		for _, r := range unicast.NLRI {
			// glog.Infof("><SB> Unicast Prefix: %+v", r.Prefix)
			n.uipv4.Add(r.Prefix, int(r.Length), peer, attr)
		}
	}

	if l3vpn, err := mpreach.GetNLRIL3VPN(); err == nil {
		glog.V(5).Info("Message with bgp mp_reach nlri vpnv4")
		// fmt.Printf("Peer: %s Attributes: %+v L3VPN Routes: %+v", peer, attr, l3vpn.NLRI)
		for _, r := range l3vpn.NLRI {
			// glog.Infof("><SB> VPN Prefix: %+v", r.Prefix)
			n.vpnv4.Add(r.Prefix, int(r.Length), peer, attr)
		}
	}
}

func (n *nlri) processMPUnReach(peer string, attr *bgp.BaseAttributes, mpunreach bgp.MPNLRI) {
	if _, err := mpunreach.GetNLRIUnicast(); err == nil {
		// fmt.Printf("Peer: %s Attributes: %+v Unicast Routes: %+v", peer, attr, unicast.NLRI)
	}
	if _, err := mpunreach.GetNLRIL3VPN(); err == nil {
		// fmt.Printf("Peer: %s Attributes: %+v L3VPN Routes: %+v", peer, attr, l3vpn.NLRI)
	}
}

func (n *nlri) GetAll(t RouteType) []string {
	switch t {
	case UnicastIPv4:
		return n.uipv4.GetAll()
	case UnicastIPv6:
		return n.uipv6.GetAll()
	case VPNv4:
		return n.vpnv4.GetAll()
	case VPNv6:
		return n.vpnv6.GetAll()
	}

	return nil
}

// NewClassifierNLRI return a new instance of a NLRI Classifier
func NewClassifierNLRI() NLRI {
	return &nlri{
		uipv4: radix.NewTree(),
		uipv6: radix.NewTree(),
		vpnv4: radix.NewTree(),
		vpnv6: radix.NewTree(),
	}
}
