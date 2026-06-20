//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/neo532/gofr"
	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gofr_layout/internal/server"
	"github.com/neo532/gokit/logger"
)

func initApp(
	context.Context,
	*config.Config,
	logger.Logger,
) (*gofr.App, func(), error) {
	panic(wire.Build(
		server.NewApi,
		newApp,
	))
}
