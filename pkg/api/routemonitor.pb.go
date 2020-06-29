// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routemonitor.proto

package routemonitor

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Prefix struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	MaskLength           uint32   `protobuf:"varint,2,opt,name=mask_length,json=maskLength,proto3" json:"mask_length,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Prefix) Reset()         { *m = Prefix{} }
func (m *Prefix) String() string { return proto.CompactTextString(m) }
func (*Prefix) ProtoMessage()    {}
func (*Prefix) Descriptor() ([]byte, []int) {
	return fileDescriptor_6ddae02f15385ea6, []int{0}
}

func (m *Prefix) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Prefix.Unmarshal(m, b)
}
func (m *Prefix) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Prefix.Marshal(b, m, deterministic)
}
func (m *Prefix) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Prefix.Merge(m, src)
}
func (m *Prefix) XXX_Size() int {
	return xxx_messageInfo_Prefix.Size(m)
}
func (m *Prefix) XXX_DiscardUnknown() {
	xxx_messageInfo_Prefix.DiscardUnknown(m)
}

var xxx_messageInfo_Prefix proto.InternalMessageInfo

func (m *Prefix) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Prefix) GetMaskLength() uint32 {
	if m != nil {
		return m.MaskLength
	}
	return 0
}

type PrefixList struct {
	PrefixList           []*Prefix `protobuf:"bytes,1,rep,name=prefix_list,json=prefixList,proto3" json:"prefix_list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PrefixList) Reset()         { *m = PrefixList{} }
func (m *PrefixList) String() string { return proto.CompactTextString(m) }
func (*PrefixList) ProtoMessage()    {}
func (*PrefixList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6ddae02f15385ea6, []int{1}
}

func (m *PrefixList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrefixList.Unmarshal(m, b)
}
func (m *PrefixList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrefixList.Marshal(b, m, deterministic)
}
func (m *PrefixList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrefixList.Merge(m, src)
}
func (m *PrefixList) XXX_Size() int {
	return xxx_messageInfo_PrefixList.Size(m)
}
func (m *PrefixList) XXX_DiscardUnknown() {
	xxx_messageInfo_PrefixList.DiscardUnknown(m)
}

var xxx_messageInfo_PrefixList proto.InternalMessageInfo

func (m *PrefixList) GetPrefixList() []*Prefix {
	if m != nil {
		return m.PrefixList
	}
	return nil
}

// MonitorRequest carries a map of prefixes by the type, supported types are:
// Unicast IPv4, Unicast IPv6, VPVv4, VPVv6
type MonitorRequest struct {
	PrefixList           map[int32]*PrefixList `protobuf:"bytes,1,rep,name=prefix_list,json=prefixList,proto3" json:"prefix_list,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *MonitorRequest) Reset()         { *m = MonitorRequest{} }
func (m *MonitorRequest) String() string { return proto.CompactTextString(m) }
func (*MonitorRequest) ProtoMessage()    {}
func (*MonitorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6ddae02f15385ea6, []int{2}
}

func (m *MonitorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorRequest.Unmarshal(m, b)
}
func (m *MonitorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorRequest.Marshal(b, m, deterministic)
}
func (m *MonitorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorRequest.Merge(m, src)
}
func (m *MonitorRequest) XXX_Size() int {
	return xxx_messageInfo_MonitorRequest.Size(m)
}
func (m *MonitorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorRequest proto.InternalMessageInfo

func (m *MonitorRequest) GetPrefixList() map[int32]*PrefixList {
	if m != nil {
		return m.PrefixList
	}
	return nil
}

// MonitorResponse carries a map of prefixes by the event,
// List of prefixes changed, list of prefixes deleted.
type MonitorResponse struct {
	PrefixList           map[int32]*PrefixList `protobuf:"bytes,1,rep,name=prefix_list,json=prefixList,proto3" json:"prefix_list,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *MonitorResponse) Reset()         { *m = MonitorResponse{} }
func (m *MonitorResponse) String() string { return proto.CompactTextString(m) }
func (*MonitorResponse) ProtoMessage()    {}
func (*MonitorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6ddae02f15385ea6, []int{3}
}

func (m *MonitorResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MonitorResponse.Unmarshal(m, b)
}
func (m *MonitorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MonitorResponse.Marshal(b, m, deterministic)
}
func (m *MonitorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MonitorResponse.Merge(m, src)
}
func (m *MonitorResponse) XXX_Size() int {
	return xxx_messageInfo_MonitorResponse.Size(m)
}
func (m *MonitorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MonitorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MonitorResponse proto.InternalMessageInfo

func (m *MonitorResponse) GetPrefixList() map[int32]*PrefixList {
	if m != nil {
		return m.PrefixList
	}
	return nil
}

func init() {
	proto.RegisterType((*Prefix)(nil), "Prefix")
	proto.RegisterType((*PrefixList)(nil), "PrefixList")
	proto.RegisterType((*MonitorRequest)(nil), "MonitorRequest")
	proto.RegisterMapType((map[int32]*PrefixList)(nil), "MonitorRequest.PrefixListEntry")
	proto.RegisterType((*MonitorResponse)(nil), "MonitorResponse")
	proto.RegisterMapType((map[int32]*PrefixList)(nil), "MonitorResponse.PrefixListEntry")
}

func init() { proto.RegisterFile("routemonitor.proto", fileDescriptor_6ddae02f15385ea6) }

var fileDescriptor_6ddae02f15385ea6 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x52, 0x4f, 0x4b, 0xc3, 0x30,
	0x14, 0x27, 0x1b, 0x5b, 0xe1, 0x65, 0xda, 0xf1, 0x4e, 0x65, 0x97, 0xd5, 0x9e, 0x7a, 0x0a, 0x52,
	0x41, 0xc4, 0x83, 0x28, 0xe2, 0x45, 0x26, 0x48, 0xbe, 0xc0, 0xa8, 0xec, 0xa9, 0x65, 0x5d, 0x53,
	0x93, 0x54, 0xdc, 0x97, 0x11, 0x3f, 0xaa, 0x34, 0x75, 0x73, 0x86, 0x9e, 0xbd, 0xbd, 0x3f, 0xbf,
	0xfc, 0xfe, 0x24, 0x01, 0xd4, 0xaa, 0xb1, 0xb4, 0x51, 0x55, 0x61, 0x95, 0x16, 0xb5, 0x56, 0x56,
	0x25, 0xb7, 0x30, 0x7e, 0xd4, 0xf4, 0x5c, 0x7c, 0x60, 0x04, 0x41, 0xbe, 0x5a, 0x69, 0x32, 0x26,
	0x62, 0x31, 0x4b, 0x27, 0x72, 0xd7, 0xe2, 0x1c, 0xf8, 0x26, 0x37, 0xeb, 0x65, 0x49, 0xd5, 0x8b,
	0x7d, 0x8d, 0x06, 0x31, 0x4b, 0x8f, 0x24, 0xb4, 0xa3, 0x85, 0x9b, 0x24, 0xe7, 0x00, 0x1d, 0xc9,
	0xa2, 0x30, 0x16, 0x53, 0xe0, 0xb5, 0xeb, 0x96, 0x65, 0x61, 0x6c, 0xc4, 0xe2, 0x61, 0xca, 0xb3,
	0x40, 0x74, 0x08, 0x09, 0xf5, 0x1e, 0x99, 0x7c, 0x32, 0x38, 0x7e, 0xe8, 0xec, 0x48, 0x7a, 0x6b,
	0xc8, 0x58, 0xbc, 0xee, 0x3b, 0x3c, 0x17, 0x7f, 0x51, 0xe2, 0x57, 0xed, 0xae, 0xb2, 0x7a, 0x7b,
	0x48, 0x3a, 0xbb, 0x87, 0xd0, 0x5b, 0xe3, 0x14, 0x86, 0x6b, 0xda, 0xba, 0x58, 0x23, 0xd9, 0x96,
	0x78, 0x02, 0xa3, 0xf7, 0xbc, 0x6c, 0xc8, 0x85, 0xe1, 0x19, 0x3f, 0x60, 0x94, 0xdd, 0xe6, 0x72,
	0x70, 0xc1, 0x92, 0x2f, 0x06, 0xe1, 0x5e, 0xda, 0xd4, 0xaa, 0x32, 0x84, 0x37, 0x7d, 0x0e, 0x63,
	0xe1, 0xc1, 0xfe, 0xcb, 0x62, 0x76, 0x05, 0x13, 0xd9, 0x3e, 0xeb, 0x8f, 0x3e, 0x0a, 0x08, 0x76,
	0x65, 0xe8, 0x5d, 0xdb, 0x6c, 0xea, 0xbb, 0x3c, 0x65, 0x4f, 0x63, 0xf7, 0x0f, 0xce, 0xbe, 0x03,
	0x00, 0x00, 0xff, 0xff, 0xf7, 0x9f, 0x0c, 0xe9, 0x1d, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RouteMonitorClient is the client API for RouteMonitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RouteMonitorClient interface {
	Monitor(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (RouteMonitor_MonitorClient, error)
}

type routeMonitorClient struct {
	cc *grpc.ClientConn
}

func NewRouteMonitorClient(cc *grpc.ClientConn) RouteMonitorClient {
	return &routeMonitorClient{cc}
}

func (c *routeMonitorClient) Monitor(ctx context.Context, in *MonitorRequest, opts ...grpc.CallOption) (RouteMonitor_MonitorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RouteMonitor_serviceDesc.Streams[0], "/RouteMonitor/Monitor", opts...)
	if err != nil {
		return nil, err
	}
	x := &routeMonitorMonitorClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RouteMonitor_MonitorClient interface {
	Recv() (*MonitorResponse, error)
	grpc.ClientStream
}

type routeMonitorMonitorClient struct {
	grpc.ClientStream
}

func (x *routeMonitorMonitorClient) Recv() (*MonitorResponse, error) {
	m := new(MonitorResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouteMonitorServer is the server API for RouteMonitor service.
type RouteMonitorServer interface {
	Monitor(*MonitorRequest, RouteMonitor_MonitorServer) error
}

// UnimplementedRouteMonitorServer can be embedded to have forward compatible implementations.
type UnimplementedRouteMonitorServer struct {
}

func (*UnimplementedRouteMonitorServer) Monitor(req *MonitorRequest, srv RouteMonitor_MonitorServer) error {
	return status.Errorf(codes.Unimplemented, "method Monitor not implemented")
}

func RegisterRouteMonitorServer(s *grpc.Server, srv RouteMonitorServer) {
	s.RegisterService(&_RouteMonitor_serviceDesc, srv)
}

func _RouteMonitor_Monitor_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MonitorRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RouteMonitorServer).Monitor(m, &routeMonitorMonitorServer{stream})
}

type RouteMonitor_MonitorServer interface {
	Send(*MonitorResponse) error
	grpc.ServerStream
}

type routeMonitorMonitorServer struct {
	grpc.ServerStream
}

func (x *routeMonitorMonitorServer) Send(m *MonitorResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _RouteMonitor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RouteMonitor",
	HandlerType: (*RouteMonitorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Monitor",
			Handler:       _RouteMonitor_Monitor_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "routemonitor.proto",
}
