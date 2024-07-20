package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gyu-oj-backend/app/user/cmd/api/internal/config"
	"gyu-oj-backend/app/user/cmd/api/internal/handler"
	"gyu-oj-backend/app/user/cmd/api/internal/svc"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/result"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf,
		rest.WithCustomCors(nil, nil, constant.AllOrigins),
		rest.WithUnauthorizedCallback(result.JwtUnauthorizedResult))

	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
