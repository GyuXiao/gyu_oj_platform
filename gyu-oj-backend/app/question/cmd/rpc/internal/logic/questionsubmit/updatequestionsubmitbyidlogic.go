package questionsubmitlogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/protobuf/encoding/protojson"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateQuestionSubmitByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateQuestionSubmitByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionSubmitByIdLogic {
	return &UpdateQuestionSubmitByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateQuestionSubmitByIdLogic) UpdateQuestionSubmitById(in *pb.QuestionSubmitUpdateReq) (*pb.QuestionSubmitUpdateResp, error) {
	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "UpdateQuestionSubmitById 的请求参数 id: "+in.Id)
	}

	questionSubmit := &entity.QuestionSubmit{}
	l.fixExtraFields(questionSubmit, in)
	_, err = do.QuestionSubmit.Where(do.QuestionSubmit.ID.Eq(int64(id))).Updates(&questionSubmit)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UpdateQuestionSubmitError)
	}

	return &pb.QuestionSubmitUpdateResp{UpdateOK: true}, nil
}

func (l *UpdateQuestionSubmitByIdLogic) fixExtraFields(questionSubmit *entity.QuestionSubmit, in *pb.QuestionSubmitUpdateReq) {
	if in.JudgeInfo != nil {
		judgeInfo, err := protojson.MarshalOptions{
			EmitUnpopulated:   true,
			EmitDefaultValues: true,
		}.Marshal(in.JudgeInfo)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		questionSubmit.JudgeInfo = string(judgeInfo)
	}
	if in.Status != 0 {
		questionSubmit.Status = in.Status
	}
}
