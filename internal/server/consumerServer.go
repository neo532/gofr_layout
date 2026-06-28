package server

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
	cg "github.com/neo532/gokit/queue/kafka/consumergroup"
	"github.com/neo532/gokit/util/slicex"
)

func connectConsumer(c context.Context,
	maxSlowtime float64,
	name string,
	group string,
	topics []string,
	addrs []string,
	log logger.Logger,
	fn func(c context.Context, message []byte) (err error),
) (*cg.ConsumerGroup, error) {
	return cg.NewGroup(
		name,
		addrs,
		group,
		cg.WithLogger(log, c),
		cg.WithTopics(topics...),
		cg.WithSlowLog(time.Duration(maxSlowtime)*time.Second),
		cg.WithAutoCommit(true),
		cg.WithBalanceStrategy(sarama.NewBalanceStrategySticky()),
		cg.WithContext(c),
		cg.WithMiddleware(),
		cg.WithHandler(func(ctx context.Context, message []byte) (err error) {
			return fn(ctx, message)
		}),
	)
}

func pkgConsumerUnit(
	c context.Context,
	cfg config.DataConsumerConfCfg,
	log logger.Logger,
	fn func(c context.Context, message []byte) (err error)) (cum queue.Consumer, err error) {

	cs := make([]queue.Consumer, 0, 3)
	var s queue.Consumer
	if s, err = connectConsumer(c,
		cfg.MaxSlowtime.Load().(float64),
		cfg.ConsumerDefault.Name.Load().(string),
		cfg.ConsumerDefault.Group.Load().(string),
		slicex.OfType[string](cfg.ConsumerDefault.Topics.Load().([]any)),
		slicex.OfType[string](cfg.ConsumerDefault.Addrs.Load().([]any)),
		log,
		fn,
	); err != nil {
		return
	}
	cs = append(cs, s)
	if _, ok := cfg.ConsumerGray.Addrs.Load().([]string); ok {
		if s, err = connectConsumer(c,
			cfg.MaxSlowtime.Load().(float64),
			cfg.ConsumerGray.Name.Load().(string),
			cfg.ConsumerGray.Group.Load().(string),
			slicex.OfType[string](cfg.ConsumerGray.Topics.Load().([]any)),
			slicex.OfType[string](cfg.ConsumerGray.Addrs.Load().([]any)),
			log,
			fn,
		); err != nil {
			return
		}
		cs = append(cs, s)
	}
	if _, ok := cfg.ConsumerShadow.Addrs.Load().([]string); ok {
		if s, err = connectConsumer(c,
			cfg.MaxSlowtime.Load().(float64),
			cfg.ConsumerShadow.Name.Load().(string),
			cfg.ConsumerShadow.Group.Load().(string),
			slicex.OfType[string](cfg.ConsumerShadow.Topics.Load().([]any)),
			slicex.OfType[string](cfg.ConsumerShadow.Addrs.Load().([]any)),
			log,
			fn,
		); err != nil {
			return
		}
		cs = append(cs, s)
	}
	cum = queue.NewConsumers(cs...)
	return
}

func NewConsumer(c context.Context, cfg *config.Config, log logger.Logger, cis []*ConsumerItem) (csm queue.Consumer, err error) {

	csms := make([]queue.Consumer, 0, 1)
	for _, v := range cis {
		var cs queue.Consumer
		if cs, err = pkgConsumerUnit(c, v.cfg, log, v.fn); err != nil {
			return
		}
		csms = append(csms, cs)
	}

	csm = queue.NewConsumers(csms...)
	return
}
