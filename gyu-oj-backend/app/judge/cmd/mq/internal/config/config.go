package config

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	JudgeRpcConf zrpc.RpcClientConf
	ListenerConf rabbitmq.RabbitListenerConf
}
