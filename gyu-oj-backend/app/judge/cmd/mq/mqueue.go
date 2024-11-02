package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"gyu-oj-backend/app/judge/cmd/mq/internal/config"
	"gyu-oj-backend/app/judge/cmd/mq/internal/listen"
)

var configFile = flag.String("f", "etc/queue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("Starting RabbitMq server at %v \n", c.ListenerConf.Port)
	err := InitRabbitMq(c)
	if err != nil {
		return
	}

	listener := listen.NewListenerService(c, context.Background())
	serviceGroup := service.NewServiceGroup()
	serviceGroup.Add(listener)
	defer serviceGroup.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()
}
