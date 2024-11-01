package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	RabbitMq struct {
		Host     string
		Port     int
		Username string
		Password string
	}
	UserRpcConf     zrpc.RpcClientConf
	QuestionRpcConf zrpc.RpcClientConf
	JudgeRpcConf    zrpc.RpcClientConf
}
