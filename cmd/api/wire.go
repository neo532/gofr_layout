//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/neo532/gofr"
	"github.com/neo532/gofr_layout/cmd"
	"github.com/neo532/gofr_layout/internal/server"
)

func initApp() (*gofr.App, func(), error) {
	panic(wire.Build(
		cmd.ProviderSet,

		server.NewApi,

		newApp,
	))
}
