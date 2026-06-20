package main

import (
	"context"
	"flag"

	"github.com/neo532/gofr"
	"github.com/neo532/gofr/transport"
	"github.com/neo532/gofr_layout/cmd"
	"github.com/neo532/gofr_layout/internal/config"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs/", "config path, eg: -conf configs")
}

func newApp(
	c context.Context,
	cfg *config.Config,
	srvs ...transport.Server,
) *gofr.App {

	return gofr.New(
		gofr.ID(cfg.General.Ip.Load().(string)),
		gofr.Name(cfg.General.Name.Load().(string)),
		gofr.Version(Version),
		gofr.Metadata(map[string]string{}),
		gofr.Context(c),
		gofr.Server(srvs...),
	)
}

func main() {
	flag.Parse()

	ctx := cmd.BootContext()

	// config
	cfg, closeCfg, err := cmd.InitConfig(ctx, flagConf)
	if err != nil {
		panic(err)
	}
	defer closeCfg()
	cfg.General.Version.Store(Version)

	// log
	log, closeLog, err := cmd.InitLogger(cfg, cmd.EntryApi)
	if err != nil {
		panic(err)
	}
	defer closeLog()

	// app
	app, cleanup, err := initApp(
		ctx,
		cfg,
		log,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// pid
	// prestop: ["cat","{path}/pid","|","xargs","kill","-2"]
	if err := app.WritePID(""); err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
