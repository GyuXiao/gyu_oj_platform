package userlogic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gyu-oj-backend/app/user/models/do"
	"gyu-oj-backend/app/user/models/token"
	"gyu-oj-backend/common/tools"
	"gyu-oj-backend/common/xerr"

	"gyu-oj-backend/app/user/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// 先通过 username 查询用户是否存在
	user, err := do.User.Where(do.User.Username.Eq(in.Username)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewErrCode(xerr.SearchUserError)
	}
	// 如果用户不存在，登陆失败，返回
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerr.NewErrCode(xerr.UserNotExistError)
	}

	// 如果用户存在，再校验用户密码是否正确
	err = checkUserPassword(user.Password, in.Password)
	if err != nil {
		return nil, err
	}

	// 用户名和密码都无误且用户存在，生成 jwt token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&GenerateTokenReq{userId: uint64(user.ID)})
	if err != nil {
		return nil, err
	}
	// token 存入缓存
	// key field value 的格式如下
	// login:token:xxx {userId: xxx, userRole: xxx, username: xxx, avatarUrl: xxx}
	tokenLogic := token.NewDefaultTokenModel(l.svcCtx.RedisClient)
	err = tokenLogic.InsertToken(tokenResp.accessToken, uint64(user.ID), uint8(user.UserRole), user.Username, user.AvatarURL)
	if err != nil {
		return nil, err
	}

	// 登陆成功，返回用户 id，用户名，token，token 过期时间
	return &pb.LoginResp{
		Id:          uint64(user.ID),
		Username:    user.Username,
		AvatarUrl:   user.AvatarURL,
		UserRole:    uint64(user.UserRole),
		Token:       tokenResp.accessToken,
		TokenExpire: tokenResp.accessExpire,
	}, nil
}

// 校验用户密码

func checkUserPassword(pwd string, password string) error {
	str, err := tools.DecodeMd5(pwd)
	if err != nil {
		return xerr.NewErrCode(xerr.DecodeMd5Error)
	}
	if !tools.DecodeBcrypt(str, password) {
		return xerr.NewErrCode(xerr.UserPasswordError)
	}
	return nil
}
