package questionlogic

import (
	"context"
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
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "QuestionDeleteReq 的请求参数 id: "+in.Id)
	}

	// 软删除
	_, err = do.Question.Where(do.Question.ID.Eq(int64(id))).Update(do.Question.IsDelete, 1)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.DeleteQuestionError)
	}

	return &pb.QuestionDeleteResp{DeleteOK: true}, nil
}
