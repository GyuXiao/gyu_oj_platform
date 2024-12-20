package questionlogic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/tools"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"
)

type AddQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddQuestionLogic {
	return &AddQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddQuestionLogic) AddQuestion(in *pb.QuestionAddReq) (*pb.QuestionAddResp, error) {
	question := &entity.Question{
		// 必要的参数
		ID:      tools.GenId(),
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserId,
		Answer:  in.Answer,
	}
	// 补充额外的非空参数
	l.fixExtraFields(question, in)
	err := do.Question.Create(question)
	if err != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.AddQuestionError), "创建一个题目错误")
	}

	id := strconv.FormatInt(question.ID, 10)
	return &pb.QuestionAddResp{Id: id}, nil
}

func (l *AddQuestionLogic) fixExtraFields(question *entity.Question, in *pb.QuestionAddReq) {
	if len(in.Tags) > 0 {
		question.Tags = strings.Join(in.Tags, ",")
	}

	if len(in.JudgeCase) > 0 {
		judgeCases, err := json.Marshal(in.JudgeCase)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		question.JudgeCases = string(judgeCases)
	}

	if in.JudgeConfig != nil {
		judgeConfig, err := json.Marshal(in.JudgeConfig)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONMarshalError))
		}
		question.JudgeConfig = string(judgeConfig)
	}
}
