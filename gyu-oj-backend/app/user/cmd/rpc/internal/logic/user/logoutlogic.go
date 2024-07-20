package userlogic

import (
	"context"
	"gyu-oj-backend/app/user/models/token"

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
	err := tokenLogic.DeleteToken(in.AuthToken)
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResp{IsLogouted: true}, nil
}
