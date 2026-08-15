//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package test

import (
	"github.com/google/wire"
	"github.com/neo532/gofr_layout/cmd"
	"github.com/neo532/gofr_layout/internal/biz"
	"github.com/neo532/gofr_layout/internal/connect"
	"github.com/neo532/gofr_layout/internal/data"
	"github.com/neo532/gofr_layout/internal/service/api"
)

func UserApi() (*api.UserApi, func(), error) {
	panic(wire.Build(
		cmd.ProviderSet,

		connect.ProviderSet,

		api.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
	))
}
