package main

import (
	"flag"
	"fmt"
	"gyu-oj-backend/common/interceptor/rpcserver"

	"gyu-oj-backend/app/question/cmd/rpc/internal/config"
	questionServer "gyu-oj-backend/app/question/cmd/rpc/internal/server/question"
	questionsubmitServer "gyu-oj-backend/app/question/cmd/rpc/internal/server/questionsubmit"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/question.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterQuestionServer(grpcServer, questionServer.NewQuestionServer(ctx))
		pb.RegisterQuestionSubmitServer(grpcServer, questionsubmitServer.NewQuestionSubmitServer(ctx))

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
