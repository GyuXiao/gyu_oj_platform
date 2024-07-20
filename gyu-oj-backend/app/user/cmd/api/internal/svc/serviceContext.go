package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/user/cmd/api/internal/config"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
