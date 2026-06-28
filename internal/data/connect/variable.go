package connect

import (
	"context"

	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gokit/database"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/logger"
	"github.com/neo532/gokit/queue"
)

type (
	DatabaseUser    *orm.Orms
	RedisLock       *redis.Rediss
	ProducerUser    queue.Producer
	TransactionUser database.Transaction
)

type WireFieldsOfDatabaseUserSet struct {
	DB DatabaseUser
	TX TransactionUser
}

func NewDatabaseUser(c context.Context, cfg *config.Config, log logger.Logger) (*WireFieldsOfDatabaseUserSet, func(), error) {
	dbs := newDatabase(c, &cfg.Data.DatabaseUser.DatabaseConf, log)
	return &WireFieldsOfDatabaseUserSet{DB: dbs, TX: dbs.Transaction}, dbs.Close(), dbs.Error()
}

func NewRedisLock(c context.Context, cfg *config.Config, log logger.Logger) (RedisLock, func(), error) {
	rdbs := newRedis(c, &cfg.Data.RedisLock.RedisConf, log)
	return rdbs, rdbs.Close(), rdbs.Error()
}

func NewProducerDefault(c context.Context, cfg *config.Config, log logger.Logger) (ProducerUser, func(), error) {
	pdcs := newProducer(c, &cfg.Data.ProducerUser.ProducerConf, log)
	return pdcs, pdcs.Close(), pdcs.Error()
}
