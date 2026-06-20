//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package test

import (
	"github.com/google/wire"
	"github.com/neo532/gofr_layout/cmd"
	"github.com/neo532/gofr_layout/internal/data"
	"github.com/neo532/gofr_layout/internal/data/base"
	"github.com/neo532/gofr_layout/internal/domain"
	"github.com/neo532/gofr_layout/internal/service/api"
)

func UserApiService() (*api.UserApiService, func(), error) {
	panic(wire.Build(
		cmd.ProviderSet,

		base.ProviderSet,

		api.ProviderSet,
		data.ProviderSet,
		domain.ProviderSet,
	))
}
