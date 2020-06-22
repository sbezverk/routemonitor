package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"net/http"
	_ "net/http/pprof"

	"github.com/golang/glog"
	"github.com/sbezverk/routemonitor/pkg/bmplistener"
)

var (
	srcPort  int
	perfPort int
)

func init() {
	flag.IntVar(&srcPort, "source-port", 5000, "port exposed to outside")
	flag.IntVar(&perfPort, "performance-port", 56767, "port used for performance debugging")
}

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt}
)

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func main() {
	flag.Parse()
	_ = flag.Set("logtostderr", "true")
	// Initializing Kafka publisher
	// other publishers sutisfying pub.Publisher interface can be used.
	go func() {
		glog.Info(http.ListenAndServe(fmt.Sprintf(":%d", perfPort), nil))
	}()
	l, err := bmplistener.NewBMPListener(srcPort, nil)
	if err != nil {
		glog.Errorf("failed to start listener with error: %+v", err)
		os.Exit(1)
	}
	l.Start()
	stopCh := setupSignalHandler()
	<-stopCh
	l.Stop()
	os.Exit(0)
}
