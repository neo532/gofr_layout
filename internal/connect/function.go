package connect

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/IBM/sarama"
	"github.com/neo532/gofr_layout/internal/config"
	kitLog "github.com/neo532/gofr_layout/kit/log"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
	cg "github.com/neo532/gokit/queue/kafka/consumergroup"
	"github.com/neo532/gokit/queue/kafka/producer"
	"github.com/neo532/gokit/util/timex"
	"gorm.io/driver/mysql"
)

func newDatabase(c context.Context, cfg *config.DataDatabaseConfCfg, log logger.Logger) (dbs *orm.Orms) {
	connet := func(c context.Context,
		d *config.DataDatabaseConfCfg, name, dsn string,
		log logger.Logger) *orm.Orm {
		return orm.New(
			name,
			mysql.Open(dsn),
			orm.WithTablePrefix(d.TablePrefix.Get()),
			orm.WithConnMaxLifetime(timex.Num2Duration(d.ConnMaxLifetime.Get(), time.Second)),
			orm.WithMaxIdleConns(int32(d.MaxIdleConns.Get())),
			orm.WithMaxOpenConns(int32(d.MaxOpenConns.Get())),
			orm.WithLogger(log),
			orm.WithSingularTable(),
			orm.WithContext(c),
			orm.WithSlowLog(timex.Num2Duration(d.MaxSlowtime.Get(), time.Second)),
			orm.WithGormProcessor(func(db *gorm.DB) {
				db.Callback().Query().After("gorm:query").Register("gokit:ignore_not_found", func(d *gorm.DB) {
					if d.Error != nil && d.RowsAffected == 0 {
						d.Error = nil
					}
				})
			}),
		)
	}
	dbs = orm.News(
		orm.WithMaster(connet(c, cfg, cfg.DatabaseMaster.Name.Get(), cfg.DatabaseMaster.Dsn.Get(), log)),
		orm.WithSlave(connet(c, cfg, cfg.DatabaseSlave.Name.Get(), cfg.DatabaseSlave.Dsn.Get(), log)),
	)
	return
}

func newRedis(c context.Context, cfg *config.DataRedisConfCfg, log logger.Logger) (rdbs *redis.Rediss) {
	connnet := func(c context.Context,
		d *config.DataRedisConfCfg,
		name, addr, password string, db int64,
		log logger.ILogger) *redis.Redis {
		return redis.New(
			name,
			addr,
			redis.WithPassword(password),
			redis.WithSlowTime(timex.Num2Duration(d.MaxSlowtime.Get(), time.Second)),
			redis.WithDb(int32(db)),
			redis.WithLogger(log),
			redis.WithContext(c),
		)
	}
	rdbs = redis.News(
		redis.WithDefault(
			connnet(c, cfg,
				cfg.RedisDefault.Name.Get(),
				cfg.RedisDefault.Addr.Get(),
				cfg.RedisDefault.Password.Get(),
				cfg.RedisDefault.Db.Get(),
				log,
			),
		),
	)
	return
}

func newProducer(c context.Context, cfg *config.DataProducerConfCfg, log logger.Logger) (pds queue.Producer) {
	connect := func(c context.Context,
		name string, addrs []string, topic string,
		log logger.Logger) queue.Producer {
		return producer.New(
			name,
			addrs,
			producer.WithTopic(topic),
			producer.WithReturnSucesses(true),
			producer.WithPartitioner(sarama.NewHashPartitioner),
			producer.WithLogger(log, c),
			producer.WithRequiredAcks(sarama.WaitForAll),
			producer.WithAsync(true),
			producer.WithMiddleware(kitLog.ProducerHeaderSet),
		)
	}
	opts := []queue.ProducerOption{
		queue.WithProducer(connect(
			c,
			cfg.ProducerDefault.Name.Get(),
			cfg.ProducerDefault.Addrs.Get(),
			cfg.ProducerDefault.Topic.Get(),
			log,
		)),
		queue.WithProducerGray(connect(
			c,
			cfg.ProducerGray.Name.Get(),
			cfg.ProducerGray.Addrs.Get(),
			cfg.ProducerGray.Topic.Get(),
			log,
		)),
		queue.WithProducerShadow(connect(
			c,
			cfg.ProducerShadow.Name.Get(),
			cfg.ProducerShadow.Addrs.Get(),
			cfg.ProducerShadow.Topic.Get(),
			log,
		)),
	}
	return queue.NewProducers(opts...)
}

func PkgConsumerUnit(
	c context.Context,
	cfg config.DataConsumerConfCfg,
	log logger.Logger,
	fn func(c context.Context, message []byte) (err error)) (cum queue.Consumer, err error) {

	connectConsumer := func(c context.Context,
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
			cg.WithSlowLog(timex.Num2Duration(maxSlowtime, time.Second)),
			cg.WithAutoCommit(true),
			cg.WithBalanceStrategy(sarama.NewBalanceStrategySticky()),
			cg.WithContext(c),
			cg.WithMiddleware(),
			cg.WithHandler(func(ctx context.Context, message []byte) (err error) {
				return fn(ctx, message)
			}),
		)
	}

	cs := make([]queue.Consumer, 0, 3)
	var s queue.Consumer
	if s, err = connectConsumer(c,
		cfg.MaxSlowtime.Get(),
		cfg.ConsumerDefault.Name.Get(),
		cfg.ConsumerDefault.Group.Get(),
		cfg.ConsumerDefault.Topics.Get(),
		cfg.ConsumerDefault.Addrs.Get(),
		log,
		fn,
	); err != nil {
		return
	}
	cs = append(cs, s)
	if cfg.ConsumerGray.Addrs.Get() != nil {
		if s, err = connectConsumer(c,
			cfg.MaxSlowtime.Get(),
			cfg.ConsumerGray.Name.Get(),
			cfg.ConsumerGray.Group.Get(),
			cfg.ConsumerGray.Topics.Get(),
			cfg.ConsumerGray.Addrs.Get(),
			log,
			fn,
		); err != nil {
			return
		}
		cs = append(cs, s)
	}
	if cfg.ConsumerShadow.Addrs.Get() != nil {
		if s, err = connectConsumer(c,
			cfg.MaxSlowtime.Get(),
			cfg.ConsumerShadow.Name.Get(),
			cfg.ConsumerShadow.Group.Get(),
			cfg.ConsumerShadow.Topics.Get(),
			cfg.ConsumerShadow.Addrs.Get(),
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
