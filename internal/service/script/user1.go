package script

import (
	"context"

	"github.com/neo532/gofr_layout/internal/biz"
	"github.com/neo532/gokit/logger"
)

type User1Script struct {
	bUser *biz.UserBiz
	log   logger.Logger
}

func NewUser1Script(
	bUser *biz.UserBiz,
	log logger.Logger,
) *User1Script {
	return &User1Script{
		bUser: bUser,
		log:   log,
	}
}

func (s *User1Script) Create(c context.Context, args ...string) (err error) {
	s.log.Infof(c, "UserScript1.Create args:%v", args)
	return
}
