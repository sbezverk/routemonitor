package server

import (
	"net"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/sbezverk/routemonitor/pkg/classifier"
	pbapi "github.com/sbezverk/routemonitor/pkg/routemonitor"
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
	id := uuid.New()
	unmonitor := make([]func(), 0)
	prefixes := make(map[int]*struct {
		t  classifier.RouteType
		id string
		p  []byte
		l  int
	})
	channels := make(map[int]chan struct{})
	i := 0
	for t, routes := range req.PrefixList {
		for _, p := range routes.PrefixList {
			channels[i] = make(chan struct{})
			if err := r.c.Monitor(classifier.RouteType(t), id.String(), p.Address, int(p.MaskLength), channels[i]); err != nil {
				// In case of error, Unmonitor all previously succesful Monitor calls, and return error to the client.
				for _, f := range unmonitor {
					f()
				}
				return err
			}
			prefixes[i] = &struct {
				t  classifier.RouteType
				id string
				p  []byte
				l  int
			}{
				t:  classifier.RouteType(t),
				id: id.String(),
				p:  p.Address,
				l:  int(p.MaskLength),
			}
			unmonitor = append(unmonitor, func() {
				r.c.Unmonitor(classifier.RouteType(t), id.String(), p.Address, int(p.MaskLength))
			})
		}
	}

	for e, c := range channels {
		select {
		// TODO msg should carry, prefix Update/Dalete, prefix and its length
		case <-c:
			// TODO, when prefix gets updated/deleted, send a message
			m := prefixes[e]
			if err := srv.Send(&pbapi.MonitorResponse{
				PrefixList: map[int32]*pbapi.PrefixList{
					int32(m.t): {
						PrefixList: []*pbapi.Prefix{
							{
								Address:    m.p,
								MaskLength: uint32(m.l),
							},
						},
					},
				},
			}); err != nil {
				for _, f := range unmonitor {
					f()
				}
				return err
			}
		}
	}

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
