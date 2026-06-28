package server

import (
	"time"

	"github.com/neo532/gofr/middleware"
	"github.com/neo532/gofr/transport/grpc"
	"github.com/neo532/gofr/transport/http"
	"github.com/neo532/gofr/transport/rpcx"
	"github.com/neo532/gofr/transport/websocket"
	"github.com/neo532/gofr_layout/internal/config"
	grpcx "github.com/smallnest/rpcx/server"
	ggrpc "google.golang.org/grpc"
)

func NewHttpServer(cfg *config.Config) *http.Server {
	return http.NewServer(
		http.Address(cfg.Server.Http.Addr.Load().(string)),
		http.Timeout(time.Duration(cfg.Server.Http.Timeout.Load()*int64(time.Second))),
		http.Middleware(
			middleware.Validator(),
		),
	)
}

func NewGrpcServer(cfg *config.Config) *grpc.Server {
	return grpc.NewServer(
		grpc.Address(cfg.Server.Grpc.Addr.Load().(string)),
		grpc.Middleware(
			middleware.Validator(),
		),
		grpc.GrpcOptions(
			ggrpc.ConnectionTimeout(time.Duration(cfg.Server.Grpc.Timeout.Load()*int64(time.Second))),
		),
	)
}

func NewRpcxerver(cfg *config.Config) *rpcx.Server {
	return rpcx.NewServer(
		rpcx.Address(cfg.Server.Rpcx.Addr.Load().(string)),
		rpcx.Middleware(
			middleware.Validator(),
		),
		rpcx.RpcxOptions(
			grpcx.WithReadTimeout(time.Duration(cfg.Server.Websocket.Timeout.Load()*int64(time.Second))),
			grpcx.WithWriteTimeout(time.Duration(cfg.Server.Websocket.Timeout.Load()*int64(time.Second))),
		),
	)
}

func NewWebsocket(cfg *config.Config) *websocket.Server {
	return websocket.NewServer(
		websocket.Address(cfg.Server.Websocket.Addr.Load().(string)),
		websocket.Timeout(time.Duration(cfg.Server.Websocket.Timeout.Load()*int64(time.Second))),
		websocket.Middleware(
			middleware.Validator(),
		),
	)
}
