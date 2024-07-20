package questionsubmitlogic

import (
	"context"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/app/question/models/enums"
	"gyu-oj-backend/common/tools"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type DoQuestionSubmitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDoQuestionSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoQuestionSubmitLogic {
	return &DoQuestionSubmitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DoQuestionSubmitLogic) DoQuestionSubmit(in *pb.QuestionSubmitAddReq) (*pb.QuestionSubmitAddResp, error) {
	questionId, err := strconv.Atoi(in.QuestionId)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "QuestionSubmitAddReq 的请求参数 questionId: "+in.QuestionId)
	}

	questionSubmit := &entity.QuestionSubmit{
		ID:         tools.GenId(),
		Language:   in.Language,
		SubmitCode: in.SubmitCode,
		JudgeInfo:  "",
		Status:     enums.WAITING,
		QuestionID: int64(questionId),
		UserID:     in.UserId,
	}
	err = do.QuestionSubmit.Create(questionSubmit)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.CreateQuestionSubmitError)
	}

	id := strconv.FormatInt(questionSubmit.ID, 10)
	return &pb.QuestionSubmitAddResp{Id: id}, nil
}
