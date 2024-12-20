// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2
// Source: question.proto

package questionsubmit

import (
	"context"

	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	JudgeCase                    = pb.JudgeCase
	JudgeConfig                  = pb.JudgeConfig
	JudgeInfo                    = pb.JudgeInfo
	PageReq                      = pb.PageReq
	QuestionAddReq               = pb.QuestionAddReq
	QuestionAddResp              = pb.QuestionAddResp
	QuestionDeleteReq            = pb.QuestionDeleteReq
	QuestionDeleteResp           = pb.QuestionDeleteResp
	QuestionGetByIdReq           = pb.QuestionGetByIdReq
	QuestionGetByIdResp          = pb.QuestionGetByIdResp
	QuestionListByPageReq        = pb.QuestionListByPageReq
	QuestionListByPageResp       = pb.QuestionListByPageResp
	QuestionSubmitAddReq         = pb.QuestionSubmitAddReq
	QuestionSubmitAddResp        = pb.QuestionSubmitAddResp
	QuestionSubmitListByPageReq  = pb.QuestionSubmitListByPageReq
	QuestionSubmitListByPageResp = pb.QuestionSubmitListByPageResp
	QuestionSubmitQueryByIdReq   = pb.QuestionSubmitQueryByIdReq
	QuestionSubmitQueryByIdResp  = pb.QuestionSubmitQueryByIdResp
	QuestionSubmitUpdateReq      = pb.QuestionSubmitUpdateReq
	QuestionSubmitUpdateResp     = pb.QuestionSubmitUpdateResp
	QuestionSubmitVO             = pb.QuestionSubmitVO
	QuestionUpdateReq            = pb.QuestionUpdateReq
	QuestionUpdateResp           = pb.QuestionUpdateResp
	QuestionVO                   = pb.QuestionVO

	QuestionSubmit interface {
		DoQuestionSubmit(ctx context.Context, in *QuestionSubmitAddReq, opts ...grpc.CallOption) (*QuestionSubmitAddResp, error)
		QueryQuestionSubmit(ctx context.Context, in *QuestionSubmitListByPageReq, opts ...grpc.CallOption) (*QuestionSubmitListByPageResp, error)
		QueryQuestionSubmitById(ctx context.Context, in *QuestionSubmitQueryByIdReq, opts ...grpc.CallOption) (*QuestionSubmitQueryByIdResp, error)
		UpdateQuestionSubmitById(ctx context.Context, in *QuestionSubmitUpdateReq, opts ...grpc.CallOption) (*QuestionSubmitUpdateResp, error)
	}

	defaultQuestionSubmit struct {
		cli zrpc.Client
	}
)

func NewQuestionSubmit(cli zrpc.Client) QuestionSubmit {
	return &defaultQuestionSubmit{
		cli: cli,
	}
}

func (m *defaultQuestionSubmit) DoQuestionSubmit(ctx context.Context, in *QuestionSubmitAddReq, opts ...grpc.CallOption) (*QuestionSubmitAddResp, error) {
	client := pb.NewQuestionSubmitClient(m.cli.Conn())
	return client.DoQuestionSubmit(ctx, in, opts...)
}

func (m *defaultQuestionSubmit) QueryQuestionSubmit(ctx context.Context, in *QuestionSubmitListByPageReq, opts ...grpc.CallOption) (*QuestionSubmitListByPageResp, error) {
	client := pb.NewQuestionSubmitClient(m.cli.Conn())
	return client.QueryQuestionSubmit(ctx, in, opts...)
}

func (m *defaultQuestionSubmit) QueryQuestionSubmitById(ctx context.Context, in *QuestionSubmitQueryByIdReq, opts ...grpc.CallOption) (*QuestionSubmitQueryByIdResp, error) {
	client := pb.NewQuestionSubmitClient(m.cli.Conn())
	return client.QueryQuestionSubmitById(ctx, in, opts...)
}

func (m *defaultQuestionSubmit) UpdateQuestionSubmitById(ctx context.Context, in *QuestionSubmitUpdateReq, opts ...grpc.CallOption) (*QuestionSubmitUpdateResp, error) {
	client := pb.NewQuestionSubmitClient(m.cli.Conn())
	return client.UpdateQuestionSubmitById(ctx, in, opts...)
}
