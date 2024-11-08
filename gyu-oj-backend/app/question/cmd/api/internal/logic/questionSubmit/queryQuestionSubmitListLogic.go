package questionSubmit

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/xerr"
	"strings"

	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryQuestionSubmitListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryQuestionSubmitListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryQuestionSubmitListLogic {
	return &QueryQuestionSubmitListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryQuestionSubmitListLogic) QueryQuestionSubmitList(req *types.QueryQuestionSubmitReq) (*types.QueryQuestionSubmitResp, error) {
	// 1,用户必须登陆后才能查看题目提交记录
	token := strings.Split(req.Authorization, " ")[1]
	currentUser, _ := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if currentUser == nil {
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	// 2,调用 rpc 模块的查询题目提交记录
	resp, err := l.svcCtx.QuestionSubmitRpc.QueryQuestionSubmit(l.ctx, &questionsubmit.QuestionSubmitListByPageReq{
		PageReq: &question.PageReq{
			Current:   req.Current,
			PageSize:  req.PageSize,
			SortField: req.SortField,
			SortOrder: req.SortOrder,
		},
		Language:   req.Language,
		Status:     req.Status,
		QuestionId: req.QuestionId,
		UserId:     req.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	var questionSubmitList []types.QuestionSubmitVO
	if len(resp.QuestionSubmitVOList) > 0 {
		err = copier.Copy(&questionSubmitList, resp.QuestionSubmitVOList)
		if err != nil {
			logc.Infof(l.ctx, "questionSubmitVOList 数据转换错误: %v\n", err)
		}
	}

	return &types.QueryQuestionSubmitResp{
		QuestionSubmitList: questionSubmitList,
		TotalNum:           resp.TotalNum,
	}, nil
}
