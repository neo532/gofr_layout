package script

import (
	"context"

	"github.com/neo532/gofr_layout/internal/biz"
	"github.com/neo532/gokit/logger"
)

type UserScript struct {
	bUser *biz.UserBiz
	log   logger.Logger
}

func NewUserScript(
	bUser *biz.UserBiz,
	log logger.Logger,
) *UserScript {
	return &UserScript{
		bUser: bUser,
		log:   log,
	}
}

func (s *UserScript) Create(c context.Context, args ...string) (err error) {
	s.log.Infof(c, "UserScript.Create args:%v", args)
	return
}
