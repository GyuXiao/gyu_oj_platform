package questionlogic

import (
	"context"
	"encoding/json"
	"errors"
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
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "QuestionGetById 的请求参数 id: "+in.Id)
	}

	question, err := do.Question.Where(do.Question.ID.Eq(int64(id))).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewErrCode(xerr.SearchQuestionByIdError)
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
		err = json.Unmarshal([]byte(question.Answer), &questionVO.Answers)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONUnmarshalError))
		}
	}

	if question.JudgeConfig != "" {
		err = json.Unmarshal([]byte(question.JudgeConfig), &questionVO.JudgeConfig)
		if err != nil {
			logc.Infof(l.ctx, xerr.GetMsgByCode(xerr.JSONUnmarshalError))
		}
	}
}
