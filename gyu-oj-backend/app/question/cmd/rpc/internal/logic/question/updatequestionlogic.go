package questionlogic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionLogic {
	return &UpdateQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateQuestionLogic) UpdateQuestion(in *pb.QuestionUpdateReq) (*pb.QuestionUpdateResp, error) {
	question := &entity.Question{}
	l.fixExtraFields(question, in)

	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "QuestionUpdateReq 的请求参数 id: "+in.Id)
	}

	_, err = do.Question.Where(do.Question.ID.Eq(int64(id))).Updates(&question)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UpdateQuestionError)
	}

	return &pb.QuestionUpdateResp{UpdateOK: true}, nil
}

func (l *UpdateQuestionLogic) fixExtraFields(question *entity.Question, in *pb.QuestionUpdateReq) {
	if len(in.Tags) > 0 {
		tags := strings.Join(in.Tags, ",")
		question.Tags = tags
	}

	if in.Title != "" {
		question.Title = in.Title
	}

	if in.Content != "" {
		question.Content = in.Content
	}

	if len(in.Answers) > 0 {
		answers, err := json.Marshal(in.Answers)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		question.Answer = string(answers)
	}

	if len(in.JudgeCases) > 0 {
		judgeCases, err := json.Marshal(in.JudgeCases)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		question.JudgeCases = string(judgeCases)
	}

	if in.JudgeConfig.TimeLimit > 0 || in.JudgeConfig.MemoryLimit > 0 {
		judgeConfig, err := json.Marshal(in.JudgeConfig)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		question.JudgeConfig = string(judgeConfig)
	}

}
