package script

import (
	"context"

	"github.com/neo532/gofr_layout/internal/domain"
	"github.com/neo532/gokit/logger"
)

type UserScript struct {
	dm  *domain.UserDomain
	log logger.Logger
}

func NewUserScript(
	dm *domain.UserDomain,
	log logger.Logger,
) *UserScript {
	return &UserScript{
		dm:  dm,
		log: log,
	}
}

func (s *UserScript) Create(c context.Context, args ...string) (err error) {
	s.log.Infof(c, "UserScript.Create args:%v", args)
	return
}
