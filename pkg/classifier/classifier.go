package classifier

import (
	"fmt"
	"net"

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

	fmt.Printf("%s peer hash: %s Update: %+v\n", peer, msg.PeerHeader.GetPeerHash(), *u)
}

// NewClassifierNLRI return a new instance of a NLRI Classifier
func NewClassifierNLRI() NLRI {
	return &nlri{}
}
