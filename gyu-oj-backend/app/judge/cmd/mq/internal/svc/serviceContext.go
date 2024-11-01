package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/judge/cmd/mq/internal/config"
	"gyu-oj-backend/app/judge/cmd/rpc/judge"
)

type ServiceContext struct {
	Config   config.Config
	JudgeRpc judge.Judge
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		JudgeRpc: judge.NewJudge(zrpc.MustNewClient(c.JudgeRpcConf)),
	}
}
