package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/judge/cmd/rpc/judge"
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
	JudgeRpc          judge.Judge
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserRpc:           user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
		QuestionRpc:       question.NewQuestion(zrpc.MustNewClient(c.QuestionRpcConf)),
		QuestionSubmitRpc: questionsubmit.NewQuestionSubmit(zrpc.MustNewClient(c.QuestionRpcConf)),
		JudgeRpc:          judge.NewJudge(zrpc.MustNewClient(c.JudgeRpcConf)),
	}
}
