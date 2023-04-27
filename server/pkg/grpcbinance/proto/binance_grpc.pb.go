// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: binance.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BinanceService_GetUSDTPrices_FullMethodName       = "/binance.BinanceService/GetUSDTPrices"
	BinanceService_Get24HChangePercent_FullMethodName = "/binance.BinanceService/Get24hChangePercent"
)

// BinanceServiceClient is the client API for BinanceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BinanceServiceClient interface {
	GetUSDTPrices(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*USDTPricesResponse, error)
	Get24HChangePercent(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ChangePercentResponse, error)
}

type binanceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBinanceServiceClient(cc grpc.ClientConnInterface) BinanceServiceClient {
	return &binanceServiceClient{cc}
}

func (c *binanceServiceClient) GetUSDTPrices(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*USDTPricesResponse, error) {
	out := new(USDTPricesResponse)
	err := c.cc.Invoke(ctx, BinanceService_GetUSDTPrices_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binanceServiceClient) Get24HChangePercent(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ChangePercentResponse, error) {
	out := new(ChangePercentResponse)
	err := c.cc.Invoke(ctx, BinanceService_Get24HChangePercent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BinanceServiceServer is the server API for BinanceService service.
// All implementations must embed UnimplementedBinanceServiceServer
// for forward compatibility
type BinanceServiceServer interface {
	GetUSDTPrices(context.Context, *Empty) (*USDTPricesResponse, error)
	Get24HChangePercent(context.Context, *Empty) (*ChangePercentResponse, error)
	mustEmbedUnimplementedBinanceServiceServer()
}

// UnimplementedBinanceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBinanceServiceServer struct {
}

func (UnimplementedBinanceServiceServer) GetUSDTPrices(context.Context, *Empty) (*USDTPricesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUSDTPrices not implemented")
}
func (UnimplementedBinanceServiceServer) Get24HChangePercent(context.Context, *Empty) (*ChangePercentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get24HChangePercent not implemented")
}
func (UnimplementedBinanceServiceServer) mustEmbedUnimplementedBinanceServiceServer() {}

// UnsafeBinanceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BinanceServiceServer will
// result in compilation errors.
type UnsafeBinanceServiceServer interface {
	mustEmbedUnimplementedBinanceServiceServer()
}

func RegisterBinanceServiceServer(s grpc.ServiceRegistrar, srv BinanceServiceServer) {
	s.RegisterService(&BinanceService_ServiceDesc, srv)
}

func _BinanceService_GetUSDTPrices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinanceServiceServer).GetUSDTPrices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinanceService_GetUSDTPrices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinanceServiceServer).GetUSDTPrices(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BinanceService_Get24HChangePercent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinanceServiceServer).Get24HChangePercent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinanceService_Get24HChangePercent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinanceServiceServer).Get24HChangePercent(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// BinanceService_ServiceDesc is the grpc.ServiceDesc for BinanceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BinanceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "binance.BinanceService",
	HandlerType: (*BinanceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUSDTPrices",
			Handler:    _BinanceService_GetUSDTPrices_Handler,
		},
		{
			MethodName: "Get24hChangePercent",
			Handler:    _BinanceService_Get24HChangePercent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "binance.proto",
}