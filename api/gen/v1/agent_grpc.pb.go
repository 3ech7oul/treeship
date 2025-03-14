// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.6
// source: v1/agent.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	AgentService_MessageRoute_FullMethodName = "/server.AgentService/MessageRoute"
)

// AgentServiceClient is the client API for AgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentServiceClient interface {
	MessageRoute(ctx context.Context, opts ...grpc.CallOption) (AgentService_MessageRouteClient, error)
}

type agentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentServiceClient(cc grpc.ClientConnInterface) AgentServiceClient {
	return &agentServiceClient{cc}
}

func (c *agentServiceClient) MessageRoute(ctx context.Context, opts ...grpc.CallOption) (AgentService_MessageRouteClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &AgentService_ServiceDesc.Streams[0], AgentService_MessageRoute_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &agentServiceMessageRouteClient{ClientStream: stream}
	return x, nil
}

type AgentService_MessageRouteClient interface {
	Send(*MessageRequest) error
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type agentServiceMessageRouteClient struct {
	grpc.ClientStream
}

func (x *agentServiceMessageRouteClient) Send(m *MessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentServiceMessageRouteClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentServiceServer is the server API for AgentService service.
// All implementations must embed UnimplementedAgentServiceServer
// for forward compatibility
type AgentServiceServer interface {
	MessageRoute(AgentService_MessageRouteServer) error
	mustEmbedUnimplementedAgentServiceServer()
}

// UnimplementedAgentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServiceServer struct {
}

func (UnimplementedAgentServiceServer) MessageRoute(AgentService_MessageRouteServer) error {
	return status.Errorf(codes.Unimplemented, "method MessageRoute not implemented")
}
func (UnimplementedAgentServiceServer) mustEmbedUnimplementedAgentServiceServer() {}

// UnsafeAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServiceServer will
// result in compilation errors.
type UnsafeAgentServiceServer interface {
	mustEmbedUnimplementedAgentServiceServer()
}

func RegisterAgentServiceServer(s grpc.ServiceRegistrar, srv AgentServiceServer) {
	s.RegisterService(&AgentService_ServiceDesc, srv)
}

func _AgentService_MessageRoute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServiceServer).MessageRoute(&agentServiceMessageRouteServer{ServerStream: stream})
}

type AgentService_MessageRouteServer interface {
	Send(*MessageResponse) error
	Recv() (*MessageRequest, error)
	grpc.ServerStream
}

type agentServiceMessageRouteServer struct {
	grpc.ServerStream
}

func (x *agentServiceMessageRouteServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentServiceMessageRouteServer) Recv() (*MessageRequest, error) {
	m := new(MessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AgentService_ServiceDesc is the grpc.ServiceDesc for AgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.AgentService",
	HandlerType: (*AgentServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MessageRoute",
			Handler:       _AgentService_MessageRoute_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/agent.proto",
}
