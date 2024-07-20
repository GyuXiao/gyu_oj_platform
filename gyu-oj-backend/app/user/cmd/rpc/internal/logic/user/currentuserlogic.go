package userlogic

import (
	"context"
	"gyu-oj-backend/app/user/models/token"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"strconv"

	"gyu-oj-backend/app/user/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentUserLogic {
	return &CurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CurrentUserLogic) CurrentUser(in *pb.CurrentUserReq) (*pb.CurrentUserResp, error) {
	// 校验从 jwt token 解析出来的 userId 是否和缓存中的 userId 一致
	// 1，从 token 中解析出 userId1
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	claims, err := generateTokenLogic.ParseTokenByKey(in.AuthToken, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}
	userId1 := claims[constant.KeyJwtUserId].(float64)

	// 2，从 根据 token 从 redis 中拿到 userId2
	tokenLogic := token.NewDefaultTokenModel(l.svcCtx.RedisClient)
	result, err := tokenLogic.CheckTokenExist(in.AuthToken)
	if err != nil {
		return nil, err
	}
	userIdStr, userRoleStr, username, avatarUrl := result[0], result[1], result[2], result[3]
	userId2, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, err
	}
	userRole, err := strconv.Atoi(userRoleStr)
	if err != nil {
		return nil, err
	}

	// 3，判断两者是否相同
	if uint64(userId1) != uint64(userId2) {
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	// 4，校验成功后，刷新 token
	tokenLogic.RefreshToken(in.AuthToken)

	return &pb.CurrentUserResp{
		Id:          uint64(userId1),
		Username:    username,
		AvatarUrl:   avatarUrl,
		UserRole:    uint64(userRole),
		Token:       in.AuthToken,
		TokenExpire: int64(constant.TokenExpireTime.Seconds()),
	}, nil
}
