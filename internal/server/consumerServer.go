package server

import (
	"context"

	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gofr_layout/internal/connect"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
)

func NewConsumer(c context.Context, cfg *config.Config, log logger.Logger, cis []*ConsumerItem) (csm queue.Consumer, err error) {

	csms := make([]queue.Consumer, 0, 1)
	for _, v := range cis {
		var cs queue.Consumer
		if cs, err = connect.PkgConsumerUnit(c, v.cfg, log, v.fn); err != nil {
			return
		}
		csms = append(csms, cs)
	}

	csm = queue.NewConsumers(csms...)
	return
}
