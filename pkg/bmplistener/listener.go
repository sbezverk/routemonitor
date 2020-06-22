package bmplistener

import (
	"fmt"
	"io"
	"net"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

// Listener defines methods to manage BMP Server
type Listener interface {
	Start()
	Stop()
}

type listener struct {
	sourcePort int
	incoming   net.Listener
	stop       chan struct{}
	classifier func(bmpMsg bmp.Message)
}

func (l *listener) Start() {
	// Starting bmp server server
	glog.Infof("Starting bmp listener on %s", l.incoming.Addr().String())
	go l.listen()
}

func (l *listener) Stop() {
	glog.Infof("Stopping bmp listener\n")
	close(l.stop)
}

func (l *listener) listen() {
	for {
		client, err := l.incoming.Accept()
		if err != nil {
			glog.Errorf("fail to accept client connection with error: %+v", err)
			continue
		}
		glog.V(5).Infof("client %+v accepted, calling bmp message processing worker", client.RemoteAddr())
		go l.peerWorker(client)
	}
}

func (l *listener) peerWorker(client net.Conn) {
	defer client.Close()

	parserQueue := make(chan []byte)
	parsStop := make(chan struct{})
	// Starting parser per client with dedicated work queue
	go l.parser(parserQueue, parsStop)
	defer func() {
		glog.V(5).Infof("all done with client %+v", client.RemoteAddr())
		close(parsStop)
	}()

	for {
		headerMsg := make([]byte, bmp.CommonHeaderLength)
		if _, err := io.ReadAtLeast(client, headerMsg, bmp.CommonHeaderLength); err != nil {
			glog.Errorf("fail to read from client %+v with error: %+v", client.RemoteAddr(), err)
			return
		}
		// Recovering common header first
		header, err := bmp.UnmarshalCommonHeader(headerMsg[:bmp.CommonHeaderLength])
		if err != nil {
			glog.Errorf("fail to recover BMP message Common Header with error: %+v", err)
			continue
		}
		// Allocating space for the message body
		msg := make([]byte, int(header.MessageLength)-bmp.CommonHeaderLength)
		if _, err := io.ReadFull(client, msg); err != nil {
			glog.Errorf("fail to read from client %+v with error: %+v", client.RemoteAddr(), err)
			return
		}

		fullMsg := make([]byte, int(header.MessageLength))
		copy(fullMsg, headerMsg)
		copy(fullMsg[bmp.CommonHeaderLength:], msg)
		parserQueue <- fullMsg
	}
}

// NewBMPListener instantiates a new instance of BMP listener
func NewBMPListener(sPort int, classifier func(bmp.Message)) (Listener, error) {
	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", sPort))
	if err != nil {
		glog.Errorf("fail to setup listener on port %d with error: %+v", sPort, err)
		return nil, err
	}
	l := listener{
		stop:       make(chan struct{}),
		sourcePort: sPort,
		incoming:   incoming,
		classifier: classifier,
	}

	return &l, nil
}
