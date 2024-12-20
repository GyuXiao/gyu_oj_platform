package userlogic

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gyu-oj-backend/app/user/models/do"
	"gyu-oj-backend/app/user/models/entity"
	"gyu-oj-backend/common/tools"
	"gyu-oj-backend/common/xerr"

	"gyu-oj-backend/app/user/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var SearchUserByUsernameError = xerr.NewErrCode(xerr.SearchUserError)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 先通过 username 查询用户是否存在
	user, err := do.User.Where(do.User.Username.Eq(in.Username)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(SearchUserByUsernameError, "通过 username 搜索 user 发生错误, err: %v", err)
	}
	// 如果存在，返回用户已经存在，注册失败
	if user != nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.UserExistError), "该 username 已经存在")
	}

	// 用户第一次注册，调用 createUser 创建用户
	// 处于数据安全考虑，用户密码存入数据库前先做加密处理
	pwd, pwdErr := encodeUserPassword(in.Password)
	if pwdErr != nil {
		return nil, errors.Wrapf(pwdErr, "加密用户密码时发生错误, 此时加密后的密码: %s", pwd)
	}
	u := &entity.User{
		Username: in.Username,
		Password: pwd,
	}
	err = do.User.Create(u)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.CreateUserError), "创建 user 发生错误, err: %v", err)
	}

	return &pb.RegisterResp{Username: in.Username}, nil
}

// 用户密码加密

func encodeUserPassword(pwd string) (string, error) {
	hashStr, err := tools.EncodeBcrypt(pwd)
	if err != nil {
		return "", xerr.NewErrCode(xerr.EncryptionError)
	}
	return tools.EncodeMd5([]byte(hashStr)), nil
}
