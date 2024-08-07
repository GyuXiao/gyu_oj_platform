package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	QuestionRpcConf zrpc.RpcClientConf
	CodeSandbox     struct {
		Type string
	}
}
