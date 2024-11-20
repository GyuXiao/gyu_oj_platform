package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/config"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
)

type ServiceContext struct {
	Config            config.Config
	QuestionRpc       question.Question
	QuestionSubmitRpc questionsubmit.QuestionSubmit
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		QuestionRpc:       question.NewQuestion(zrpc.MustNewClient(c.QuestionRpcConf)),
		QuestionSubmitRpc: questionsubmit.NewQuestionSubmit(zrpc.MustNewClient(c.QuestionRpcConf)),
	}
}
