package main

import (
	"context"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/judge/cmd/mq/internal/config"
)

func InitRabbitMq(c config.Config) error {
	ctx := context.Background()
	conf := rabbitmq.RabbitConf{
		Host:     c.ListenerConf.Host,
		Port:     c.ListenerConf.Port,
		Username: c.ListenerConf.Username,
		Password: c.ListenerConf.Password,
	}
	admin := rabbitmq.MustNewAdmin(conf)
	exchangeConf := rabbitmq.ExchangeConf{
		ExchangeName: "oj_exchange",
		Type:         "direct",
		Durable:      true,
		AutoDelete:   false,
		Internal:     false,
		NoWait:       false,
	}

	err := admin.DeclareExchange(exchangeConf, nil)
	if err != nil {
		logc.Infof(ctx, "RabbitMq 申明一个交换机错误: %v", err)
		return err
	}

	queueConf := rabbitmq.QueueConf{
		Name:       "oj_queue",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
	err = admin.DeclareQueue(queueConf, nil)
	if err != nil {
		logc.Infof(ctx, "RabbitMq 申明一个队列错误: %v", err)
		return err
	}

	err = admin.Bind(queueConf.Name, "oj_routingKey", exchangeConf.ExchangeName, false, nil)
	if err != nil {
		logc.Infof(ctx, "RabbitMq 队列绑定交换机错误: %v", err)
		return err
	}
	return nil
}
