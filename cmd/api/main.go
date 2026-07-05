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
)

func init() {
	flag.StringVar(&cmd.ConfigPath, "conf", cmd.ConfigPath, "config path, eg: -conf configs")
	cmd.Entry = "api"
	cmd.Version = Version
}

func newApp(
	c context.Context,
	cfg *config.Config,
	srvs ...transport.Server,
) *gofr.App {

	config.Cfg.General.Version.Store(Version)

	return gofr.New(
		gofr.ID(cfg.General.Ip.Load().(string)),
		gofr.Name(cfg.General.Name.Load().(string)),
		gofr.Version(Version),
		gofr.Metadata(map[string]string{}),
		gofr.Context(c),
		gofr.EnableUpgrader(),
		gofr.Server(srvs...),
	)
}

func main() {
	flag.Parse()

	// app
	app, cleanup, err := initApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
