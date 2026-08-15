package consumer

import (
	"context"

	"github.com/neo532/gofr_layout/internal/biz"
	"github.com/neo532/gokit/logger"
)

type UserConsumer struct {
	bUser *biz.UserBiz
	log   logger.Logger
}

func NewUserConsumer(
	bUser *biz.UserBiz,
	log logger.Logger,
) *UserConsumer {
	return &UserConsumer{
		bUser: bUser,
		log:   log,
	}
}

func (s *UserConsumer) Create(c context.Context, message []byte) (err error) {
	s.log.Infof(c, "UserConsumer.Create, msg,%v", string(message))
	return
}
