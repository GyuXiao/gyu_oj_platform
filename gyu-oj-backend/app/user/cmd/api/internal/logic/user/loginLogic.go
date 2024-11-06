package user

import (
	"context"
	"github.com/pkg/errors"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"regexp"

	"gyu-oj-backend/app/user/cmd/api/internal/svc"
	"gyu-oj-backend/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 校验参数
	if req.Username == constant.BlankString || req.Password == constant.BlankString || len(req.Username) < constant.UsernameMinLen || len(req.Password) < constant.PasswordMinLen {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.RequestParamError), "用户名或密码错误")
	}
	_, err = regexp.MatchString(constant.PatternStr, req.Username)
	if err != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.RequestParamError), "用户名称包含非法字符")
	}

	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	return &types.LoginResp{
		Id:          loginResp.Id,
		Username:    loginResp.Username,
		AvatarUrl:   loginResp.AvatarUrl,
		UserRole:    uint8(loginResp.UserRole),
		Token:       loginResp.Token,
		TokenExpire: loginResp.TokenExpire,
	}, nil
}
