package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/user/cmd/api/internal/svc"
	"gyu-oj-backend/app/user/cmd/api/internal/types"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/ctxdata"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"

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

	userId := ctxdata.GetUserIdFromCtx(l.ctx)
	if userId < 0 {
		logc.Infof(l.ctx, "从 context 中获得的 userId: %v", userId)
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	currentResp, err := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{
		AuthToken: token,
		UserId:    strconv.FormatInt(userId, 10),
	})
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
