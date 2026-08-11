package server

import (
	"github.com/dqixuan/stock_info/configs"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *configs.GrpcServer) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Addr != "" {
		opts = append(opts, grpc.Address(c.Addr))
	}
	if c.Timeout != 0 {
		opts = append(opts, grpc.Timeout(c.Timeout))
	}
	srv := grpc.NewServer(opts...)
	return srv
}
