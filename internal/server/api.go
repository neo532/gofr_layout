package server

import (
	"github.com/neo532/gofr/transport"
	"github.com/neo532/gofr/transport/grpc"
	"github.com/neo532/gofr/transport/http"
	"github.com/neo532/gofr/transport/rpcx"
	"github.com/neo532/gofr_layout/internal/service/api"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gofr/transport/websocket"
)

// NewApi registers services onto existing HTTP and gRPC servers.
// It returns the servers as []transport.Server for wire injection into newApp.
func NewApi(hs *http.Server, gs *grpc.Server, rx *rpcx.Server, ws *websocket.Server,
	UserApi *api.UserApi,
	User1Api *api.User1Api,
) []transport.Server {

	pb.RegisterHTTPServer(hs,
		UserApi,
		User1Api,
	)
	pb.RegisterGRPCServer(gs,
		UserApi,
		User1Api,
	)
	pb.RegisterRPCXServer(rx,
		UserApi,
		User1Api,
	)
	pb.RegisterWebsocketServer(ws,
		UserApi,
		User1Api,
	)

	return []transport.Server{hs, gs, rx, ws}
}
