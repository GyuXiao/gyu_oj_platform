package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/user/cmd/api/internal/config"
	"gyu-oj-backend/app/user/cmd/api/internal/middleware"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
)

type ServiceContext struct {
	Config            config.Config
	UserRpc           user.UserZrpcClient
	JwtAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserRpc:           user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
		JwtAuthMiddleware: middleware.NewJwtAuthMiddleware(c.Auth.AccessSecret).Handle,
	}
}
