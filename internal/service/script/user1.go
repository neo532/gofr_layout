package script

import (
	"context"

	"github.com/neo532/gofr_layout/internal/domain"
	"github.com/neo532/gokit/logger"
)

type User1Script struct {
	dm  *domain.UserDomain
	log logger.Logger
}

func NewUser1Script(
	dm *domain.UserDomain,
	log logger.Logger,
) *User1Script {
	return &User1Script{
		dm:  dm,
		log: log,
	}
}

func (s *User1Script) Create(c context.Context, args ...string) (err error) {
	s.log.Infof(c, "User1Script.Create args:%v", args)
	return
}
