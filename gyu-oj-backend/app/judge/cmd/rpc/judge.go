package main

import (
	"flag"
	"fmt"
	"gyu-oj-backend/common/interceptor/rpcserver"

	"gyu-oj-backend/app/judge/cmd/rpc/internal/config"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/server"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/judge/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/judge.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterJudgeServer(grpcServer, server.NewJudgeServer(ctx))

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
