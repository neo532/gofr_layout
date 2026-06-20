package server

import (
	"time"

	"github.com/neo532/gofr/middleware"
	"github.com/neo532/gofr/transport"
	"github.com/neo532/gofr/transport/http"
	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gofr_layout/internal/service/api"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

func NewApi(cfg *config.Config) (svcs []transport.Server) {

	svcs = make([]transport.Server, 0, 2)

	opts := []http.ServerOption{
		http.Address(cfg.Server.Http.Addr.Load().(string)),
		http.Timeout(time.Duration(cfg.Server.Http.Timeout.Load()) * time.Second),
		http.Middleware(
			middleware.Validator(),
		),
	}
	httpSrv := http.NewServer(opts...)
	// httpSrv.GET("/hello/:Name", func(ctx http.Context) error {
	// 	return ctx.String(200, "custom: "+ctx.PathValue("Name"))
	// })
	pb.RegisterHTTPServer(httpSrv, &api.UserApiService{}, &api.UserApi1Service{})
	svcs = append(svcs, httpSrv)

	// var opts = []khttp.ServerOption{
	// 	khttp.Middleware(
	// 		server.SetEnv(bs.General.Env),
	// 		server.SetEntry(middleware.EntryApi),
	// 		tracing.Server(),
	// 		recovery.Recovery(),
	// 		log.Server(logging),
	// 		validate.Validator(),
	// 	),
	// 	khttp.ResponseEncoder(http.ResponseEncoder),
	// 	khttp.ErrorEncoder(http.ErrorEncoder),
	// }
	// if bs.Server.Http.Network != "" {
	// 	opts = append(opts, khttp.Network(bs.Server.Http.Network))
	// }
	// if bs.Server.Http.Addr != "" {
	// 	opts = append(opts, khttp.Address(bs.Server.Http.Addr))
	// }
	// if bs.Server.Http.Timeout != nil {
	// 	opts = append(opts, khttp.Timeout(bs.Server.Http.Timeout.AsDuration()))
	// }
	// srv := khttp.NewServer(opts...)

	return
}
