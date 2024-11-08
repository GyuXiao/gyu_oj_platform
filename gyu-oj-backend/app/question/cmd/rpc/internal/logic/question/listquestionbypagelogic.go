package questionlogic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gen"
	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"
)

type ListQuestionByPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListQuestionByPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionByPageLogic {
	return &ListQuestionByPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListQuestionByPageLogic) ListQuestionByPage(in *pb.QuestionListByPageReq) (*pb.QuestionListByPageResp, error) {
	orderField, ok := do.Question.GetFieldByName(in.PageReq.SortField)
	if !ok {
		// 默认按照 ID 排序
		orderField = do.Question.ID
	}
	// 默认是降序
	orderCon := orderField.Desc()
	if in.PageReq.SortOrder == constant.Asc {
		orderCon = orderField.Asc()
	}
	limit := in.PageReq.PageSize
	offset := (in.PageReq.Current - 1) * limit
	// 根据传入的字段是否合法来构建 WHERE 查询条件
	var whereCon []gen.Condition
	if len(in.Tags) > 0 {
		tags := l.getRegexpStrings(in.Tags)
		whereCon = append(whereCon, do.Question.Where(do.Question.Tags.Regexp(strings.Join(tags, "|"))))
	}
	if in.Title != "" {
		whereCon = append(whereCon, do.Question.Where(do.Question.Title.Like(in.Title+"%")))
	}
	// 分页查找
	questionList, totalCnt, err := do.Question.Where(whereCon...).Where(do.Question.IsDelete.Eq(0)).Order(orderCon).FindByPage(int(offset), int(limit))
	if err != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.SearchQuestionPageListError), "分页查询 question 列表错误")
	}

	questionVOList := make([]*pb.QuestionVO, len(questionList))
	l.fixList(questionList, questionVOList)

	return &pb.QuestionListByPageResp{
		QuestionVOList: questionVOList,
		TotalNum:       totalCnt,
	}, nil
}

func (l *ListQuestionByPageLogic) fixList(sourceList []*entity.Question, questionVOList []*pb.QuestionVO) {
	for i, question := range sourceList {
		questionVO := pb.QuestionVO{
			Id:          strconv.FormatInt(question.ID, 10),
			Title:       question.Title,
			Content:     question.Content,
			SubmitNum:   question.SubmitNum,
			AcceptedNum: question.AcceptedNum,
			UserId:      question.UserID,
			CreateTime:  question.CreateTime.Unix(),
			UpdateTime:  question.UpdateTime.Unix(),
		}
		l.fixExtraFields(question, &questionVO)
		questionVOList[i] = &questionVO
	}
}

func (l *ListQuestionByPageLogic) fixExtraFields(question *entity.Question, questionVO *pb.QuestionVO) {
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

func (l *ListQuestionByPageLogic) getRegexpStrings(tags []string) []string {
	var res []string
	if len(tags) == 0 {
		res = append(res, ".*")
		return res
	}
	for _, tag := range tags {
		res = append(res, "(^|,)"+tag+"(|,)")
	}
	return res
}
