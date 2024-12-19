package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/ctxdata"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"

	"gyu-oj-backend/app/user/cmd/api/internal/svc"
	"gyu-oj-backend/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.LogoutResp, error) {
	token := strings.Split(req.Authorization, " ")[1]

	userId := ctxdata.GetUserIdFromCtx(l.ctx)
	if userId < 0 {
		logc.Infof(l.ctx, "从 context 中获得的 userId: %v", userId)
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	logoutResp, err := l.svcCtx.UserRpc.Logout(l.ctx, &user.LogoutReq{
		AuthToken: token,
		UserId:    strconv.FormatInt(userId, 10),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.LogoutResp{IsLogouted: logoutResp.IsLogouted}, nil
}
