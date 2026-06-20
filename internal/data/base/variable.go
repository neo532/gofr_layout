package base

import (
	"context"

	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/logger"
)

type (
	DatabaseDefault *orm.Orms
	RedisLock       *redis.Rediss
)

func NewDatabaseDefault(c context.Context, cfg *config.Config, log logger.Logger) (DatabaseDefault, func(), error) {
	dbs := newDatabase(c, &cfg.Data.DatabaseDefault.DatabaseConf, log)
	return dbs, dbs.Close(), dbs.Error()
}

func NewRedisLock(c context.Context, cfg *config.Config, log logger.Logger) (RedisLock, func(), error) {
	rdbs := newRedis(c, &cfg.Data.RedisLock.RedisConf, log)
	return rdbs, rdbs.Close(), rdbs.Error()
}
