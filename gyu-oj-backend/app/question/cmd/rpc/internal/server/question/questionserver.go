// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2
// Source: question.proto

package server

import (
	"context"

	"gyu-oj-backend/app/question/cmd/rpc/internal/logic/question"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
)

type QuestionServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedQuestionServer
}

func NewQuestionServer(svcCtx *svc.ServiceContext) *QuestionServer {
	return &QuestionServer{
		svcCtx: svcCtx,
	}
}

func (s *QuestionServer) AddQuestion(ctx context.Context, in *pb.QuestionAddReq) (*pb.QuestionAddResp, error) {
	l := questionlogic.NewAddQuestionLogic(ctx, s.svcCtx)
	return l.AddQuestion(in)
}

func (s *QuestionServer) UpdateQuestion(ctx context.Context, in *pb.QuestionUpdateReq) (*pb.QuestionUpdateResp, error) {
	l := questionlogic.NewUpdateQuestionLogic(ctx, s.svcCtx)
	return l.UpdateQuestion(in)
}

func (s *QuestionServer) DeleteQuestion(ctx context.Context, in *pb.QuestionDeleteReq) (*pb.QuestionDeleteResp, error) {
	l := questionlogic.NewDeleteQuestionLogic(ctx, s.svcCtx)
	return l.DeleteQuestion(in)
}

func (s *QuestionServer) GetQuestionById(ctx context.Context, in *pb.QuestionGetByIdReq) (*pb.QuestionGetByIdResp, error) {
	l := questionlogic.NewGetQuestionByIdLogic(ctx, s.svcCtx)
	return l.GetQuestionById(in)
}

func (s *QuestionServer) ListQuestionByPage(ctx context.Context, in *pb.QuestionListByPageReq) (*pb.QuestionListByPageResp, error) {
	l := questionlogic.NewListQuestionByPageLogic(ctx, s.svcCtx)
	return l.ListQuestionByPage(in)
}
