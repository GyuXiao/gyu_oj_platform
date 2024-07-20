package svc

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gyu-oj-backend/app/user/cmd/rpc/internal/config"
	"gyu-oj-backend/app/user/models/do"
	"log"
	"os"
	"time"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 连接数据库
	dbLog := logger.New(log.New(os.Stdout, "user-rpc ", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
		LogLevel:      logger.Error,
	})
	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{
		Logger: dbLog,
	})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}

	logc.Info(context.Background(), "user-rpc-server connect MySQL database success")
	do.SetDefault(db)

	// 业务配置
	return &ServiceContext{
		Config: c,
		RedisClient: redis.MustNewRedis(redis.RedisConf{
			Host: c.Redis.Host,
			Type: c.Redis.Type,
			Pass: c.Redis.Pass,
		}),
	}
}
