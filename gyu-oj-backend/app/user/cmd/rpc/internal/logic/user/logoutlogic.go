package userlogic

import (
	"context"
	"github.com/pkg/errors"
	"gyu-oj-backend/app/user/models/token"
	"gyu-oj-backend/common/xerr"

	"gyu-oj-backend/app/user/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *pb.LogoutReq) (*pb.LogoutResp, error) {
	tokenLogic := token.NewDefaultTokenModel(l.svcCtx.RedisClient)
	// 从缓存中取出用户信息
	result, err := tokenLogic.CheckTokenExist(in.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "从缓存中解析 token 发生错误, err: %v, result: %v", err, result)
	}

	tokenFromCache := result[len(result)-1]
	// 判断 authToken 是否与缓存中的 token 一致
	if tokenFromCache != in.AuthToken {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.TokenInvalid), "jwtToken 和从缓存中取出来的 token 不同，可能 jwtToken 被篡改")
	}

	// 退出登陆，删除缓存中的 token
	err = tokenLogic.DeleteToken(in.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "删除 Token 时发生错误")
	}
	return &pb.LogoutResp{IsLogouted: true}, nil
}
