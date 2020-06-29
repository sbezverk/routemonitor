package server

import (
	"net"

	"github.com/golang/glog"
	"github.com/google/uuid"
	pbapi "github.com/sbezverk/routemonitor/pkg/api"
	"github.com/sbezverk/routemonitor/pkg/classifier"
	"google.golang.org/grpc"
)

// RouteMonitor defines interface to RouteMonitor gRPC server
type RouteMonitor interface {
	Start()
	Stop()
}

type routeMonitor struct {
	gSrv *grpc.Server
	conn net.Listener
	c    classifier.NLRI
}

func (r *routeMonitor) Start() {
	glog.V(3).Infof("Starting RouteMonitor's gRPC on %s", r.conn.Addr().String())

	go r.gSrv.Serve(r.conn)
}

func (r *routeMonitor) Stop() {
	glog.V(3).Infof("Stopping RouteMonitor's gRPC server...")
	// First stopping grpc server
	r.gSrv.Stop()
}

func (r *routeMonitor) Monitor(req *pbapi.MonitorRequest, srv pbapi.RouteMonitor_MonitorServer) error {
	// Generating unique id for the client request
	id = uuid.New()
	for t, routes := range req.PrefixList {
		for _, r := range routes.PrefixList {

		}
	}
	// cm, err := c.Recv()
	// if err != nil {
	// 	return err
	// }
	// if cm == nil {
	// 	return fmt.Errorf("client info is nil")
	// }
	// glog.V(5).Infof("request from client with id %s to monitor.", string(cm.Id))
	// if m := g.clientMgmt.Get(string(cm.Id)); m == nil {
	// 	g.clientMgmt.Add(string(cm.Id))
	// 	glog.V(5).Infof("adding client with id %s to the store.", string(cm.Id))
	// } else {
	// 	// TODO add better handling of such condition
	// 	glog.Warningf("duplicate monitor request, client with id: %s already in the store", string(cm.Id))
	// 	return err
	// }
	// for {
	// 	_, err := c.Recv()
	// 	if err != nil {
	// 		// Error indicates that the client is no longer functional, sending command to
	// 		// the clients manager to remove the client and exit
	// 		glog.V(5).Infof("client with id %s is no longer alive, error: %+v, deleting it from the store.", string(cm.Id), err)
	// 		c := g.clientMgmt.Get(string(cm.Id))
	// 		for _, f := range c.GetRouteCleanup() {
	// 			if err := f(); err != nil {
	// 				glog.Errorf("route cleanup encountered error: %+v", err)
	// 			}
	// 		}
	// 		g.clientMgmt.Delete(string(cm.Id))
	// 		return err
	// 	}
	// }
	return nil
}

// NewRouteMonitor return an instance of RouteMonitor interface
func NewRouteMonitor(conn net.Listener, c classifier.NLRI) RouteMonitor {
	gSrv := routeMonitor{
		conn: conn,
		gSrv: grpc.NewServer([]grpc.ServerOption{}...),
	}
	pbapi.RegisterRouteMonitorServer(gSrv.gSrv, &gSrv)

	return &gSrv
}
