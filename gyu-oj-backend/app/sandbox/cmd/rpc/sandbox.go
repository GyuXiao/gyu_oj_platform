package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/config"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/server"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"
	"gyu-oj-backend/common/interceptor/rpcserver"
)

var configFile = flag.String("f", "etc/sandbox.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterCodeSandboxServer(grpcServer, server.NewCodeSandboxServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// rpc log
	s.AddUnaryInterceptors(rpcserver.LoggerInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
