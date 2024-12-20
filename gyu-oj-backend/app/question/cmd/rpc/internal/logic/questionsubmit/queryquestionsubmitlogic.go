package questionsubmitlogic

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gen"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryQuestionSubmitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryQuestionSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryQuestionSubmitLogic {
	return &QueryQuestionSubmitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryQuestionSubmitLogic) QueryQuestionSubmit(in *pb.QuestionSubmitListByPageReq) (*pb.QuestionSubmitListByPageResp, error) {
	orderField, ok := do.QuestionSubmit.GetFieldByName(in.PageReq.SortField)
	if !ok {
		orderField = do.QuestionSubmit.ID
	}
	orderCon := orderField.Desc()
	if in.PageReq.SortOrder == constant.Asc {
		orderCon = orderField.Asc()
	}
	limit := in.PageReq.PageSize
	offset := (in.PageReq.Current - 1) * limit
	// 根据传入的字段是否合法来构建 WHERE 查询条件
	var whereCon []gen.Condition
	if in.Language != "" {
		whereCon = append(whereCon, do.QuestionSubmit.Where(do.QuestionSubmit.Language.Eq(in.Language)))
	}
	if in.QuestionId != "" {
		questionId, _ := strconv.Atoi(in.QuestionId)
		whereCon = append(whereCon, do.QuestionSubmit.Where(do.QuestionSubmit.QuestionID.Eq(int64(questionId))))
	}
	if in.UserId != 0 {
		whereCon = append(whereCon, do.QuestionSubmit.Where(do.QuestionSubmit.UserID.Eq(in.UserId)))
	}
	// 分页查询
	questionSubmitList, totalCnt, err := do.QuestionSubmit.Where(whereCon...).Order(orderCon).FindByPage(int(offset), int(limit))
	if err != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.QueryQuestionSubmitError), "分页查询 questionSubmit 错误")
	}

	questionSubmitVOList := make([]*pb.QuestionSubmitVO, totalCnt)
	l.fixList(questionSubmitList, questionSubmitVOList)

	return &pb.QuestionSubmitListByPageResp{
		QuestionSubmitVOList: questionSubmitVOList,
		TotalNum:             totalCnt,
	}, nil
}

func (l *QueryQuestionSubmitLogic) fixList(questionSubmitList []*entity.QuestionSubmit, questionSubmitVOList []*pb.QuestionSubmitVO) {
	for i, questionSubmit := range questionSubmitList {
		questionSubmitVO := pb.QuestionSubmitVO{
			Id:         strconv.FormatInt(questionSubmit.ID, 10),
			Language:   questionSubmit.Language,
			SubmitCode: questionSubmit.SubmitCode,
			Status:     questionSubmit.Status,
			QuestionId: strconv.FormatInt(questionSubmit.QuestionID, 10),
			UserId:     questionSubmit.UserID,
			CreateTime: questionSubmit.CreateTime.Unix(),
			UpdateTime: questionSubmit.UpdateTime.Unix(),
		}
		FieldsConvert(questionSubmit, &questionSubmitVO)
		questionSubmitVOList[i] = &questionSubmitVO
	}
}

func FieldsConvert(questionSubmit *entity.QuestionSubmit, questionSubmitVO *pb.QuestionSubmitVO) {
	if questionSubmit.JudgeInfo != "" {
		var tmp struct {
			Message string
			Time    string
			Memory  string
		}
		err := json.Unmarshal([]byte(questionSubmit.JudgeInfo), &tmp)
		if err != nil {
			logc.Infof(context.Background(), "jsonUnmarshal err: %v", err)
		}

		var time, memory = 0, 0
		if tmp.Time != "" {
			time, _ = strconv.Atoi(tmp.Time)
		}
		if tmp.Memory != "" {
			memory, _ = strconv.Atoi(tmp.Memory)
		}

		questionSubmitVO.JudgeInfo = &pb.JudgeInfo{
			Message: tmp.Message,
			Time:    int64(time),
			Memory:  int64(memory),
		}
	}
}
