package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/question/cmd/api/internal/config"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
)

type ServiceContext struct {
	Config            config.Config
	UserRpc           user.UserZrpcClient
	QuestionRpc       question.Question
	QuestionSubmitRpc questionsubmit.QuestionSubmit
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserRpc:           user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
		QuestionRpc:       question.NewQuestion(zrpc.MustNewClient(c.QuestionRpcConf)),
		QuestionSubmitRpc: questionsubmit.NewQuestionSubmit(zrpc.MustNewClient(c.QuestionRpcConf)),
	}
}
