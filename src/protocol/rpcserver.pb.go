
package rpcserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

const _ = proto.ProtoPackageIsVersion2

type String struct {
	Str string `protobuf:"bytes,1,opt,name=str" json:"str,omitempty"`
}

func (m *String) Reset()                    { *m = String{} }
func (m *String) String() string            { return proto.CompactTextString(m) }
func (*String) ProtoMessage()               {}
func (*String) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *String) GetStr() string {
	if m != nil {
		return m.Str
	}
	return ""
}

type ServerQuery struct {
	MessageId uint64 `protobuf:"varint,1,opt,name=message_id,json=messageId" json:"message_id,omitempty"`
	GetRam    bool   `protobuf:"varint,2,opt,name=get_ram,json=getRam" json:"get_ram,omitempty"`
	GetCpu    bool   `protobuf:"varint,3,opt,name=get_cpu,json=getCpu" json:"get_cpu,omitempty"`
	Command   string `protobuf:"bytes,4,opt,name=command" json:"command,omitempty"`
}

func (m *ServerQuery) Reset()                    { *m = ServerQuery{} }
func (m *ServerQuery) String() string            { return proto.CompactTextString(m) }
func (*ServerQuery) ProtoMessage()               {}
func (*ServerQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ServerQuery) GetMessageId() uint64 {
	if m != nil {
		return m.MessageId
	}
	return 0
}

func (m *ServerQuery) GetGetRam() bool {
	if m != nil {
		return m.GetRam
	}
	return false
}

func (m *ServerQuery) GetGetCpu() bool {
	if m != nil {
		return m.GetCpu
	}
	return false
}

func (m *ServerQuery) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

type ServerReply struct {
	Messages  []string `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	MessageId uint64   `protobuf:"varint,2,opt,name=message_id,json=messageId" json:"message_id,omitempty"`
	RamUsage  string   `protobuf:"bytes,3,opt,n.
ame=ram_usage,json=ramUsage" json:"ram_usage,omitempty"`
	CpuUsage  string   `protobuf:"bytes,4,opt,name=cpu_usage,json=cpuUsage" json:"cpu_usage,omitempty"`
}

func (m *ServerReply) Reset()                    { *m = ServerReply{} }
func (m *ServerReply) String() string            { return proto.CompactTextString(m) }
func (*ServerReply) ProtoMessage()               {}
func (*ServerReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ServerReply) GetMessages() []string {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *ServerReply) GetMessageId() uint64 {
	if m != nil {
		return m.MessageId
	}
	return 0
}

func (m *ServerReply) GetRamUsage() string {
	if m != nil {
		return m.RamUsage
	}
	return ""
}

func (m *ServerReply) GetCpuUsage() string {
	if m != nil {
		return m.CpuUsage
	}
	return ""
}

func init() {
	proto.RegisterType((*String)(nil), "String")
	proto.RegisterType((*ServerQuery)(nil), "ServerQuery")
	proto.RegisterType((*ServerReply)(nil), "ServerReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RPCServer service

type RPCServerClient interface {
	Version(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	List(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Stop(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Start(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Kill(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	InstanceStop(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Attach(ctx context.Context, opts ...grpc.CallOption) (RPCServer_AttachClient, error)
}

type rPCServerClient struct {
	cc *grpc.ClientConn
}

func NewRPCServerClient(cc *grpc.ClientConn) RPCServerClient {
	return &rPCServerClient{cc}
}

func (c *rPCServerClient) Version(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/Version", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) List(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) Stop(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) Start(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/Start", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) Kill(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/Kill", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) InstanceStop(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := grpc.Invoke(ctx, "/RPCServer/InstanceStop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCServerClient) Attach(ctx context.Context, opts ...grpc.CallOption) (RPCServer_AttachClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_RPCServer_serviceDesc.Streams[0], c.cc, "/RPCServer/Attach", opts...)
	if err != nil {
		return nil, err
	}
	x := &rPCServerAttachClient{stream}
	return x, nil
}

type RPCServer_AttachClient interface {
	Send(*ServerQuery) error
	Recv() (*ServerReply, error)
	grpc.ClientStream
}

type rPCServerAttachClient struct {
	grpc.ClientStream
}

func (x *rPCServerAttachClient) Send(m *ServerQuery) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rPCServerAttachClient) Recv() (*ServerReply, error) {
	m := new(ServerReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for RPCServer service

type RPCServerServer interface {
	Version(context.Context, *String) (*String, error)
	List(context.Context, *String) (*String, error)
	Stop(context.Context, *String) (*String, error)
	Start(context.Context, *String) (*String, error)
	Kill(context.Context, *String) (*String, error)
	InstanceStop(context.Context, *String) (*String, error)
	Attach(RPCServer_AttachServer) error
}

func RegisterRPCServerServer(s *grpc.Server, srv RPCServerServer) {
	s.RegisterService(&_RPCServer_serviceDesc, srv)
}

func _RPCServer_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).Version(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).List(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).Stop(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).Start(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_Kill_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).Kill(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/Kill",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).Kill(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_InstanceStop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCServerServer).InstanceStop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCServer/InstanceStop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCServerServer).InstanceStop(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCServer_Attach_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RPCServerServer).Attach(&rPCServerAttachServer{stream})
}

type RPCServer_AttachServer interface {
	Send(*ServerReply) error
	Recv() (*ServerQuery, error)
	grpc.ServerStream
}

type rPCServerAttachServer struct {
	grpc.ServerStream
}

func (x *rPCServerAttachServer) Send(m *ServerReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rPCServerAttachServer) Recv() (*ServerQuery, error) {
	m := new(ServerQuery)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _RPCServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RPCServer",
	HandlerType: (*RPCServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _RPCServer_Version_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RPCServer_List_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _RPCServer_Stop_Handler,
		},
		{
			MethodName: "Start",
			Handler:    _RPCServer_Start_Handler,
		},
		{
			MethodName: "Kill",
			Handler:    _RPCServer_Kill_Handler,
		},
		{
			MethodName: "InstanceStop",
			Handler:    _RPCServer_InstanceStop_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Attach",
			Handler:       _RPCServer_Attach_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "rpcserver.proto",
}

func init() { proto.RegisterFile("rpcserver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x4b, 0x33, 0x31,
	0x10, 0xc6, 0x49, 0xdb, 0x77, 0xff, 0xcc, 0x5b, 0x50, 0x72, 0x31, 0x6e, 0x11, 0x96, 0x9e, 0xf6,
	0x54, 0x44, 0x3f, 0x81, 0xf4, 0x54, 0xf4, 0xa0, 0x59, 0xf4, 0x5a, 0x62, 0x1a, 0xd6, 0x85, 0x66,
	0x37, 0x4c, 0x12, 0xb1, 0x57, 0xbf, 0xa7, 0xdf, 0x45, 0x36, 0xed, 0x16, 0x29, 0xdb, 0xdb, 0x3c,
	0xf9, 0x4d, 0xf2, 0x3c, 0x4c, 0x06, 0x2e, 0xd0, 0x48, 0xab, 0xf0, 0x53, 0xe1, 0xc2, 0x60, 0xeb,
	0xda, 0x79, 0x06, 0x51, 0xe9, 0xb0, 0x6e, 0x2a, 0x7a, 0x09, 0x63, 0xeb, 0x90, 0x91, 0x9c, 0x14,
	0x29, 0xef, 0xca, 0xf9, 0x17, 0xfc, 0x2f, 0x43, 0xef, 0x8b, 0x57, 0xb8, 0xa3, 0x37, 0x00, 0x5a,
	0x59, 0x2b, 0x2a, 0xb5, 0xae, 0x37, 0xa1, 0x6f, 0xc2, 0xd3, 0xc3, 0xc9, 0x6a, 0x43, 0xaf, 0x20,
	0xae, 0x94, 0x5b, 0xa3, 0xd0, 0x6c, 0x94, 0x93, 0x22, 0xe1, 0x51, 0xa5, 0x1c, 0x17, 0xba, 0x07,
	0xd2, 0x78, 0x36, 0x3e, 0x82, 0xa5, 0xf1, 0x94, 0x41, 0x2c, 0x5b, 0xad, 0x45, 0xb3, 0x61, 0x93,
	0xe0, 0xda, 0xcb, 0xf9, 0x37, 0xe9, 0xad, 0xb9, 0x32, 0xdb, 0x1d, 0xcd, 0x20, 0x39, 0x18, 0x59,
	0x46, 0xf2, 0x71, 0x91, 0xf2, 0xa3, 0x3e, 0x89, 0x35, 0x3a, 0x8d, 0x35, 0x83, 0x14, 0x85, 0x5e,
	0xfb, 0x4e, 0x06, 0xff, 0x94, 0x27, 0x28, 0xf4, 0x6b, 0xa7, 0x3b, 0x28, 0x8d, 0x3f, 0xc0, 0x7d,
	0x86, 0x44, 0x1a, 0x1f, 0xe0, 0xdd, 0x0f, 0x81, 0x94, 0x3f, 0x2f, 0xf7, 0x39, 0xe8, 0x0c, 0xe2,
	0x37, 0x85, 0xb6, 0x6e, 0x1b, 0x1a, 0x2f, 0xf6, 0x23, 0xcb, 0xfa, 0x82, 0x32, 0x98, 0x3c, 0xd5,
	0xd6, 0x0d, 0x93, 0xd2, 0xb5, 0x66, 0x80, 0x5c, 0xc3, 0xbf, 0xd2, 0x09, 0x3c, 0x73, 0xe9, 0xb1,
	0xde, 0x6e, 0x07, 0x48, 0x0e, 0xd3, 0x55, 0x63, 0x9d, 0x68, 0xa4, 0x3a, 0xf3, 0x6c, 0x01, 0xd1,
	0x83, 0x73, 0x42, 0x7e, 0xd0, 0xe9, 0xe2, 0xcf, 0xef, 0x65, 0xbd, 0x0a, 0x03, 0x2d, 0xc8, 0x2d,
	0x79, 0x8f, 0xc2, 0x06, 0xdc, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x5e, 0xbe, 0xb6, 0xc2, 0x14,
	0x02, 0x00, 0x00,
}
