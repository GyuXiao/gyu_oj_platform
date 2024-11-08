package question

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/common/xerr"

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
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	var list []types.QuestionVO
	if len(resp.QuestionVOList) > 0 {
		err = copier.Copy(&list, resp.QuestionVOList)
		if err != nil {
			logc.Infof(l.ctx, "questionVOList 数据转换错误: %v\n", err)
		}
		for i := range list {
			l.fixExtraFields(&list[i], resp.QuestionVOList[i])
		}
	}
	return &types.GetQuestionListResp{
		QuestionList: list,
		Total:        resp.TotalNum,
	}, nil
}

func (l *QueryQuestionListLogic) fixExtraFields(questionVO *types.QuestionVO, resp *pb.QuestionVO) {
	judgeConfigBytes, err := json.Marshal(resp.JudgeConfig)
	if err != nil {
		logc.Infof(l.ctx, "judgeConfig 对象转换为 byte 数组错误: %v\n", xerr.NewErrCode(xerr.JSONMarshalError))
	}
	if string(judgeConfigBytes) == "null" {
		judgeConfigBytes = []byte("")
	}
	questionVO.JudgeConfig = string(judgeConfigBytes)

	judgeCaseBytes, err := json.Marshal(resp.JudgeCase)
	if err != nil {
		logc.Infof(l.ctx, "judgeCase 对象转换为 byte 数组错误: %v\n", xerr.NewErrCode(xerr.JSONMarshalError))
	}
	if string(judgeCaseBytes) == "null" {
		judgeCaseBytes = []byte("")
	}
	questionVO.JudgeCase = string(judgeCaseBytes)
}
