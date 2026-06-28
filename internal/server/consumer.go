package server

import (
	"context"

	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gofr_layout/internal/service/consumer"
)

type ConsumerItem struct {
	cfg config.DataConsumerConfCfg
	fn  func(c context.Context, message []byte) (err error)
}

func NewConsumerRouter(cfg *config.Config,
	user *consumer.UserConsumer,
) []*ConsumerItem {
	return []*ConsumerItem{
		{cfg.Data.ConsumerUser.ConsumerConf, user.Create},
	}
}
