package classifier

import (
	"net"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

// NLRI defines method to classify a type of NLRI found in BMP Route Monitor message
type NLRI interface {
	Classify(bmp.Message)
}

type nlri struct {
}

func (n *nlri) Classify(msg bmp.Message) {
	if net.IP(msg.PeerHeader.PeerAddress).To4() != nil {
		glog.Infof("Peer: %s", net.IP(msg.PeerHeader.PeerAddress).To4().String())
	} else {
		glog.Infof("Peer: %s", net.IP(msg.PeerHeader.PeerAddress).To16().String())
	}
}

// NewClassifierNLRI return a new instance of a NLRI Classifier
func NewClassifierNLRI() NLRI {
	return &nlri{}
}
