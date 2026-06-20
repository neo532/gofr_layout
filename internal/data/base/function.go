package base

import (
	"context"
	"time"

	"github.com/neo532/gofr_layout/internal/config"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/logger"
	"gorm.io/driver/mysql"
)

func newDatabase(c context.Context, cfg *config.DataDatabaseConfCfg, log logger.ILogger) (dbs *orm.Orms) {
	connet := func(c context.Context, d *config.DataDatabaseConfCfg, name, dsn string, log logger.ILogger) *orm.Orm {
		return orm.New(
			name,
			mysql.Open(dsn),
			orm.WithTablePrefix(d.TablePrefix.Load().(string)),
			orm.WithConnMaxLifetime(time.Duration(d.ConnMaxLifetime.Load())*time.Second),
			orm.WithMaxIdleConns(int32(d.MaxIdleConns.Load())),
			orm.WithMaxOpenConns(int32(d.MaxOpenConns.Load())),
			orm.WithLogger(log),
			orm.WithSingularTable(),
			orm.WithContext(c),
			orm.WithSlowLog(time.Duration(d.MaxSlowtime.Load().(float64))*time.Second),
		)
	}
	dbs = orm.News(
		orm.WithRead(connet(c, cfg, cfg.DatabaseRead.Name.Load().(string), cfg.DatabaseRead.Dsn.Load().(string), log)),
		orm.WithWrite(connet(c, cfg, cfg.DatabaseWrite.Name.Load().(string), cfg.DatabaseWrite.Dsn.Load().(string), log)),
	)
	return
}

func newRedis(c context.Context, cfg *config.DataRedisConfCfg, log logger.Logger) (rdbs *redis.Rediss) {
	connnet := func(c context.Context,
		d *config.DataRedisConfCfg,
		name, addr, password string, db int64,
		logger logger.ILogger) *redis.Redis {
		return redis.New(
			name,
			addr,
			redis.WithPassword(password),
			redis.WithSlowTime(time.Duration(d.MaxSlowtime.Load().(float64))*time.Second),
			redis.WithDb(int32(db)),
			redis.WithLogger(logger),
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
