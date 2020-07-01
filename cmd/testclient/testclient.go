package main

import (
	"context"
	"flag"
	"net"
	"os"

	"github.com/golang/glog"
	"github.com/sbezverk/routemonitor/pkg/classifier"
	pbapi "github.com/sbezverk/routemonitor/pkg/routemonitor"
	"google.golang.org/grpc"
)

var (
	routeMonitor string
)

func init() {
	flag.StringVar(&routeMonitor, "gateway", "localhost:5000", "Address to access route monitor")
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")

	conn, err := grpc.DialContext(context.TODO(), routeMonitor, grpc.WithInsecure())
	if err != nil {
		glog.Errorf("failed to connect to route monitor at the address: %s with error: %+v", routeMonitor, err)
		os.Exit(1)
	}
	defer conn.Close()
	rmClient := pbapi.NewRouteMonitorClient(conn)

	ms, err := rmClient.Monitor(context.TODO(), &pbapi.MonitorRequest{
		PrefixList: map[int32]*pbapi.PrefixList{
			int32(classifier.VPNv4): {
				PrefixList: []*pbapi.Prefix{
					{
						Address:    net.ParseIP("192.168.5.1").To4(),
						MaskLength: 32,
					},
				},
			},
		},
	})
	if err != nil {
		glog.Errorf("failed to get the route monitor stream with error: %+v", err)
		os.Exit(1)
	}
	for {
		resp, err := ms.Recv()
		if err != nil {
			glog.Errorf("failed to receive route monitor message from the monitor stream with error: %+v", err)
			os.Exit(1)
		}
		pl := resp.PrefixList[int32(classifier.VPNv4)]
		glog.Infof("received prefixes:")
		for _, p := range pl.PrefixList {
			if p == nil {
				continue
			}
			glog.Infof("- %s/%d", net.IP(p.Address).To4().String(), p.MaskLength)
		}
	}
}
