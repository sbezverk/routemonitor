package classifier

import (
	"fmt"
	"net"

	"github.com/sbezverk/gobmp/pkg/base"
	"github.com/sbezverk/gobmp/pkg/bgp"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

// NLRI defines method to classify a type of NLRI found in BMP Route Monitor message
type NLRI interface {
	Classify(bmp.Message)
}

type nlri struct {
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
	fmt.Printf("Update carries rfc4271 NLRI\n")
	fmt.Printf("Peer: %s Attributes: %+v Routes: %+v", peer, attr, routes)

}

func (n *nlri) processBGPWithdraw(peer string, attr *bgp.BaseAttributes, routes []base.Route) {
	fmt.Printf("Update carries rfc4271 Withdraw routes\n")
	fmt.Printf("Peer: %s Attributes: %+vRoutes: %+v", peer, attr, routes)
}
func (n *nlri) processMPReach(peer string, attr *bgp.BaseAttributes, mpreach bgp.MPNLRI) {
	if unicast, err := mpreach.GetNLRIUnicast(); err == nil {
		fmt.Printf("Peer: %s Attributes: %+v Unicast Routes: %+v", peer, attr, unicast.NLRI)
	}
	if l3vpn, err := mpreach.GetNLRIL3VPN(); err == nil {
		fmt.Printf("Peer: %s Attributes: %+v L3VPN Routes: %+v", peer, attr, l3vpn.NLRI)
	}
}

func (n *nlri) processMPUnReach(peer string, attr *bgp.BaseAttributes, mpunreach bgp.MPNLRI) {
	if unicast, err := mpunreach.GetNLRIUnicast(); err == nil {
		fmt.Printf("Peer: %s Attributes: %+v Unicast Routes: %+v", peer, attr, unicast.NLRI)
	}
	if l3vpn, err := mpunreach.GetNLRIL3VPN(); err == nil {
		fmt.Printf("Peer: %s Attributes: %+v L3VPN Routes: %+v", peer, attr, l3vpn.NLRI)
	}
}

// NewClassifierNLRI return a new instance of a NLRI Classifier
func NewClassifierNLRI() NLRI {
	return &nlri{}
}
