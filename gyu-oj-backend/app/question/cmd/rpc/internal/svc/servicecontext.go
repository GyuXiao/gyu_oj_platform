package svc

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gyu-oj-backend/app/question/cmd/rpc/internal/config"
	"gyu-oj-backend/app/question/models/do"
	"log"
	"os"
	"time"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbLog := logger.New(log.New(os.Stdout, "question-rpc ", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
		LogLevel:      logger.Error,
	})
	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{Logger: dbLog})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}

	logc.Info(context.Background(), "Question module connect MySQL database success")
	do.SetDefault(db)

	return &ServiceContext{Config: c}
}
