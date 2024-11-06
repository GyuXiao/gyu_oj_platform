package user

import (
	"context"
	"github.com/pkg/errors"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"strings"

	"gyu-oj-backend/app/user/cmd/api/internal/svc"
	"gyu-oj-backend/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentLogic {
	return &CurrentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentLogic) Current(req *types.CurrentUserReq) (resp *types.CurrentUserResp, err error) {
	token := strings.Split(req.Authorization, " ")[1]
	currentResp, err := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.CurrentUserResp{
		Id:          currentResp.Id,
		Username:    currentResp.Username,
		AvatarUrl:   currentResp.AvatarUrl,
		UserRole:    uint8(currentResp.UserRole),
		Token:       currentResp.Token,
		TokenExpire: currentResp.TokenExpire,
	}, nil
}
