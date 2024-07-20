package question

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/common/xerr"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryQuestionLogic {
	return &QueryQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryQuestionLogic) QueryQuestion(req *types.GetQuestionReq) (*types.GetQuestionResp, error) {
	// 校验参数
	if req.Id == "" {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "GetQuestionReq 的 id 不合法")
	}

	// 调用 rpc 模块根据 id 获取 question
	resp, err := l.svcCtx.QuestionRpc.GetQuestionById(l.ctx, &question.QuestionGetByIdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	var questionVO types.QuestionVO
	if resp.QuestionVO != nil {
		err = copier.Copy(&questionVO, resp.QuestionVO)
		if err != nil {
			logc.Infof(l.ctx, "questionVO 数据转换错误: %v\n", err)
		}
	}
	return &types.GetQuestionResp{Question: questionVO}, nil
}
