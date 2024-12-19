package userlogic

import (
	"context"
	"github.com/pkg/errors"
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
	tokenLogic := token.NewDefaultTokenModel(l.svcCtx.RedisClient)
	result, err := tokenLogic.CheckTokenExist(in.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "从缓存中解析 token 发生错误, err: %v, result: %v", err, result)
	}

	userIdStr, userRoleStr, username, avatarUrl, tokenFromCache := result[0], result[1], result[2], result[3], result[4]
	// 校验 authToken 是否与缓存中的 token 一致
	if tokenFromCache != in.AuthToken {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.TokenInvalid), "authToken 和从缓存中取出来的 token 不同，可能 authToken 被篡改了")
	}

	// 验证通过，则刷新 token
	tokenLogic.RefreshToken(in.UserId)

	// 返回用户信息
	userIdFromCache, _ := strconv.Atoi(userIdStr)
	userRole, _ := strconv.Atoi(userRoleStr)
	return &pb.CurrentUserResp{
		Id:          uint64(userIdFromCache),
		Username:    username,
		AvatarUrl:   avatarUrl,
		UserRole:    uint64(userRole),
		Token:       tokenFromCache,
		TokenExpire: int64(constant.TokenExpireTime.Seconds()),
	}, nil
}
