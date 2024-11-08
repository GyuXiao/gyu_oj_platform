package question

import (
	"context"
	"github.com/pkg/errors"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"strings"

	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(req *types.DeleteQuestionReq) (*types.DeleteQuestionResp, error) {
	// 1,非管理员不能删除题目
	token := strings.Split(req.Authorization, " ")[1]
	currentUser, _ := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if currentUser != nil && currentUser.UserRole != constant.AdminRole {
		return nil, xerr.NewErrCode(xerr.UserNotAdminError)
	}

	// 2,调用 rpc 模块删除题目
	resp, err := l.svcCtx.QuestionRpc.DeleteQuestion(l.ctx, &question.QuestionDeleteReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.DeleteQuestionResp{IsDeleted: resp.DeleteOK}, nil
}
