package svc

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	DockerClient *client.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{Config: c}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logc.Infof(context.Background(), "创建 docker client 失败 err: %v", err)
	}
	if cli != nil {
		svc.DockerClient = cli
	}

	return svc
}
