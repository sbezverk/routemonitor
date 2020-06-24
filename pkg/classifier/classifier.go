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
	GetAll(RouteType) []string
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
	case VPNv4:
	case VPNv6:
	}

	return false
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
	// fmt.Printf("Update carries rfc4271 NLRI\n")
	// fmt.Printf("Peer: %s Attributes: %+v Routes: %+v", peer, attr, routes)
	for _, r := range routes {
		n.uipv4.Add(r.Prefix, int(r.Length), peer, attr)
	}
}

func (n *nlri) processBGPWithdraw(peer string, attr *bgp.BaseAttributes, routes []base.Route) {
	// fmt.Printf("Update carries rfc4271 Withdraw routes\n")
	// fmt.Printf("Peer: %s Attributes: %+vRoutes: %+v", peer, attr, routes)
}
func (n *nlri) processMPReach(peer string, attr *bgp.BaseAttributes, mpreach bgp.MPNLRI) {
	if unicast, err := mpreach.GetNLRIUnicast(); err == nil {
		// fmt.Printf("Peer: %s Attributes: %+v Unicast Routes: %+v", peer, attr, unicast.NLRI)
		for _, r := range unicast.NLRI {
			n.uipv4.Add(r.Prefix, int(r.Length), peer, attr)
		}
		all := n.vpnv4.GetAll()
		glog.Infof("All MP Reach Unicast: %+v", all)
		if len(all) != 0 {
			p, m, _ := net.ParseCIDR(all[0])
			l, _ := m.Mask.Size()
			if !n.vpnv4.Check(net.IP(p).To4(), l) {
				glog.Errorf("Check for existing prefix failed")
			} else {
				glog.Infof("Check for existing prefix succeeded")
			}
		}
	}
	if l3vpn, err := mpreach.GetNLRIL3VPN(); err == nil {
		// fmt.Printf("Peer: %s Attributes: %+v L3VPN Routes: %+v", peer, attr, l3vpn.NLRI)
		for _, r := range l3vpn.NLRI {
			n.vpnv4.Add(r.Prefix, int(r.Length), peer, attr)
		}
		all := n.vpnv4.GetAll()
		glog.Infof("All MP Reach VPNv4: %+v", all)
		if len(all) != 0 {
			p, m, _ := net.ParseCIDR(all[0])
			l, _ := m.Mask.Size()
			if !n.vpnv4.Check(net.IP(p).To4(), l) {
				glog.Errorf("Check for existing prefix failed")
			} else {
				glog.Infof("Check for existing prefix succeeded")
			}
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
