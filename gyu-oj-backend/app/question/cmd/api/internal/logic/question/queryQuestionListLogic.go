package question

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"

	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryQuestionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryQuestionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryQuestionListLogic {
	return &QueryQuestionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryQuestionListLogic) QueryQuestionList(req *types.GetQuestionListReq) (*types.GetQuestionListResp, error) {
	// 调用 rpc 模块获取 []questionVO 列表
	resp, err := l.svcCtx.QuestionRpc.ListQuestionByPage(l.ctx, &question.QuestionListByPageReq{
		PageReq: &question.PageReq{
			Current:   req.Current,
			PageSize:  req.PageSize,
			SortField: req.SortField,
			SortOrder: req.SortOrder,
		},
		Title: req.Title,
		Tags:  req.Tags,
	})
	if err != nil {
		return nil, err
	}

	var list []types.QuestionVO
	if len(resp.QuestionVOList) > 0 {
		err = copier.Copy(&list, resp.QuestionVOList)
		if err != nil {
			logc.Infof(l.ctx, "questionVOList 数据转换错误: %v\n", err)
		}
	}
	return &types.GetQuestionListResp{
		QuestionList: list,
		Total:        resp.TotalNum,
	}, nil
}
