package questionlogic

import (
	"context"
	"github.com/pkg/errors"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(in *pb.QuestionDeleteReq) (*pb.QuestionDeleteResp, error) {
	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ParamFormatError), "格式错误, questionId: %s", in.Id)
	}

	// 软删除
	_, err = do.Question.Where(do.Question.ID.Eq(int64(id))).Update(do.Question.IsDelete, 1)
	if err != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.DeleteQuestionError), "删除题目错误")
	}

	return &pb.QuestionDeleteResp{DeleteOK: true}, nil
}
