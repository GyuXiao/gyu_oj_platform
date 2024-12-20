package questionsubmitlogic

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gyu-oj-backend/app/question/models/do"
	"gyu-oj-backend/app/question/models/entity"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"gyu-oj-backend/app/question/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/question/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryQuestionSubmitByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryQuestionSubmitByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryQuestionSubmitByIdLogic {
	return &QueryQuestionSubmitByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryQuestionSubmitByIdLogic) QueryQuestionSubmitById(in *pb.QuestionSubmitQueryByIdReq) (*pb.QuestionSubmitQueryByIdResp, error) {
	id, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ParamFormatError), "格式转换错误, questionSubmitId: %s", in.Id)
	}
	questionSubmit, err := do.QuestionSubmit.Where(do.QuestionSubmit.ID.Eq(int64(id))).First()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.QuestionSubmitNotExistError), "questionSubmit 记录不存在")
	}
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.QueryQuestionSubmitByIdError), "通过 id 查询 questionSubmit 错误, id: %s", in.Id)
	}
	questionSubmitVO := pb.QuestionSubmitVO{
		// 必要参数
		Id:         strconv.FormatInt(questionSubmit.ID, 10),
		Language:   questionSubmit.Language,
		SubmitCode: questionSubmit.SubmitCode,
		Status:     questionSubmit.Status,
		QuestionId: strconv.FormatInt(questionSubmit.QuestionID, 10),
		UserId:     questionSubmit.UserID,
		CreateTime: questionSubmit.CreateTime.Unix(),
		UpdateTime: questionSubmit.UpdateTime.Unix(),
	}
	// 补充额外参数
	l.fixExtraFields(questionSubmit, &questionSubmitVO)

	return &pb.QuestionSubmitQueryByIdResp{QuestionSubmitVO: &questionSubmitVO}, nil
}

func (l *QueryQuestionSubmitByIdLogic) fixExtraFields(questionSubmit *entity.QuestionSubmit, questionSubmitVO *pb.QuestionSubmitVO) {
	if questionSubmit.JudgeInfo != "" {
		FieldsConvert(questionSubmit, questionSubmitVO)
	}
}
