package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/config"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
	"gyu-oj-backend/app/sandbox/cmd/rpc/codesandbox"
)

type ServiceContext struct {
	Config            config.Config
	QuestionRpc       question.Question
	QuestionSubmitRpc questionsubmit.QuestionSubmit
	SandboxRpc        codesandbox.CodeSandbox
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		QuestionRpc:       question.NewQuestion(zrpc.MustNewClient(c.QuestionRpcConf)),
		QuestionSubmitRpc: questionsubmit.NewQuestionSubmit(zrpc.MustNewClient(c.QuestionRpcConf)),
		SandboxRpc:        codesandbox.NewCodeSandbox(zrpc.MustNewClient(c.SandboxRpcConf)),
	}
}
