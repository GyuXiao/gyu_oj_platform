package questionlogic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"
)

type GetQuestionByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionByIdLogic {
	return &GetQuestionByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionByIdLogic) GetQuestionById(in *pb.QuestionGetByIdReq) (*pb.QuestionGetByIdResp, error) {
	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ParamFormatError), "格式错误, questionId: %s", in.Id)
	}

	question, err := do.Question.Where(do.Question.ID.Eq(int64(id))).First()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.QuestionNotExistError), "question 记录不存在")
	}
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SearchQuestionByIdError), "通过 id 查询 question 错误, id: %s", in.Id)
	}
	questionVO := pb.QuestionVO{
		// 必要参数
		Id:          strconv.FormatInt(question.ID, 10),
		Title:       question.Title,
		Content:     question.Content,
		SubmitNum:   question.SubmitNum,
		AcceptedNum: question.AcceptedNum,
		UserId:      question.UserID,
		CreateTime:  question.CreateTime.Unix(),
		UpdateTime:  question.UpdateTime.Unix(),
	}
	// 补充额外的非空参数
	l.fixExtraFields(question, &questionVO)

	return &pb.QuestionGetByIdResp{QuestionVO: &questionVO}, nil
}

func (l *GetQuestionByIdLogic) fixExtraFields(question *entity.Question, questionVO *pb.QuestionVO) {
	if question.Tags != "" {
		questionVO.Tags = strings.Split(question.Tags, ",")
	}

	var err error
	if question.Answer != "" {
		questionVO.Answer = question.Answer
	}

	if question.JudgeConfig != "" {
		err = json.Unmarshal([]byte(question.JudgeConfig), &questionVO.JudgeConfig)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONUnmarshalError))
		}
	}

	if question.JudgeCases != "" {
		err = json.Unmarshal([]byte(question.JudgeCases), &questionVO.JudgeCase)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONUnmarshalError))
		}
	}
}
