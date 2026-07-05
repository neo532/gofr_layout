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
	"github.com/neo532/gokit/util/slicex"
	"gorm.io/driver/mysql"
)

func newDatabase(c context.Context, cfg *config.DataDatabaseConfCfg, log logger.Logger) (dbs *orm.Orms) {
	connet := func(c context.Context,
		d *config.DataDatabaseConfCfg, name, dsn string,
		log logger.Logger) *orm.Orm {
		return orm.New(
			name,
			mysql.Open(dsn),
			orm.WithTablePrefix(d.TablePrefix.Load().(string)),
			orm.WithConnMaxLifetime(time.Duration(d.ConnMaxLifetime.Load()*int64(time.Second))),
			orm.WithMaxIdleConns(int32(d.MaxIdleConns.Load())),
			orm.WithMaxOpenConns(int32(d.MaxOpenConns.Load())),
			orm.WithLogger(log),
			orm.WithSingularTable(),
			orm.WithContext(c),
			orm.WithSlowLog(time.Duration(d.MaxSlowtime.Load().(float64)*float64(time.Second))),
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
		orm.WithMaster(connet(c, cfg, cfg.DatabaseMaster.Name.Load().(string), cfg.DatabaseMaster.Dsn.Load().(string), log)),
		orm.WithSlave(connet(c, cfg, cfg.DatabaseSlave.Name.Load().(string), cfg.DatabaseSlave.Dsn.Load().(string), log)),
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
			redis.WithSlowTime(time.Duration(d.MaxSlowtime.Load().(float64)*float64(time.Second))),
			redis.WithDb(int32(db)),
			redis.WithLogger(log),
			redis.WithContext(c),
		)
	}
	rdbs = redis.News(
		redis.WithDefault(
			connnet(c, cfg,
				cfg.RedisDefault.Name.Load().(string),
				cfg.RedisDefault.Addr.Load().(string),
				cfg.RedisDefault.Password.Load().(string),
				cfg.RedisDefault.Db.Load(),
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
			cfg.ProducerDefault.Name.Load().(string),
			slicex.OfType[string](cfg.ProducerDefault.Addrs.Load().([]any)),
			cfg.ProducerDefault.Topic.Load().(string),
			log,
		)),
		queue.WithProducerGray(connect(
			c,
			cfg.ProducerGray.Name.Load().(string),
			slicex.OfType[string](cfg.ProducerGray.Addrs.Load().([]any)),
			cfg.ProducerGray.Topic.Load().(string),
			log,
		)),
		queue.WithProducerShadow(connect(
			c,
			cfg.ProducerShadow.Name.Load().(string),
			slicex.OfType[string](cfg.ProducerShadow.Addrs.Load().([]any)),
			cfg.ProducerShadow.Topic.Load().(string),
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
