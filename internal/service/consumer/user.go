package consumer

import (
	"context"

	"github.com/neo532/gofr_layout/internal/domain"
	"github.com/neo532/gokit/logger"
)

type UserConsumer struct {
	dm  *domain.UserDomain
	log logger.Logger
}

func NewUserConsumer(
	dm *domain.UserDomain,
	log logger.Logger,
) *UserConsumer {
	return &UserConsumer{
		dm:  dm,
		log: log,
	}
}

func (s *UserConsumer) Create(c context.Context, message []byte) (err error) {
	s.log.Infof(c, "UserConsumer.Create, msg,%v", string(message))
	return
}
