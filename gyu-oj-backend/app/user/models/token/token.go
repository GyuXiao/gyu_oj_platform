package token

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"sync"
)

/*
* 使用 Redis 进行存放，校验，刷新，删除 token 信息
 */

type JwtTokenModel interface {
	InsertToken(string, uint64, uint8, string, string) error
	CheckTokenExist(string) ([]string, error)
	RefreshToken(string)
	DeleteToken(string) error
}

var ctx = context.Background()
var tokenModel JwtTokenModel
var tokenOnce sync.Once

type defaultTokenModel struct {
	*redis.Redis
}

func NewDefaultTokenModel(rds *redis.Redis) JwtTokenModel {
	tokenOnce.Do(func() {
		tokenModel = &defaultTokenModel{rds}
	})
	return tokenModel
}

// 向 Redis 中插入一条 token 记录，记录 UserId 和 UserRole 信息

func (rds *defaultTokenModel) InsertToken(token string, userId uint64, userRole uint8, username string, avatarUrl string) error {
	key := constant.TokenPrefixStr + token
	err := rds.PipelinedCtx(ctx, func(pipeline redis.Pipeliner) error {
		pipeline.HSet(ctx, key,
			constant.KeyUserId, userId,
			constant.KeyUserRole, userRole,
			constant.KeyUsername, username,
			constant.KeyAvatarUrl, avatarUrl)
		pipeline.Expire(ctx, key, constant.TokenExpireTime)
		_, execErr := pipeline.Exec(ctx)
		return execErr
	})
	if err != nil {
		logc.Infof(ctx, "redis insert token by userId err: %v", err)
		return xerr.NewErrCode(xerr.KeyInsertError)
	}
	return nil
}

// 判断 Redis 中是否存在对应的 token 记录

func (rds *defaultTokenModel) CheckTokenExist(token string) ([]string, error) {
	key := constant.TokenPrefixStr + token
	// result 的格式是 [userId, userRole, username, avatarUrl]
	result, err := rds.HmgetCtx(ctx, key, constant.KeyUserId, constant.KeyUserRole, constant.KeyUsername, constant.KeyAvatarUrl)
	if err != nil {
		logc.Infof(ctx, "redis HMGet key err: %v", err)
		return nil, xerr.NewErrCode(xerr.TokenGetFromCacheError)
	}
	if result[0] == "" || result[1] == "" {
		return nil, xerr.NewErrMsg("根据 token 拿到的数据为空字符串")
	}
	return result, nil
}

// 刷新 token 的过期时间

func (rds *defaultTokenModel) RefreshToken(token string) {
	_, err := rds.CheckTokenExist(token)
	if err != nil {
		logc.Info(ctx, err)
	}
	key := constant.TokenPrefixStr + token
	err = rds.ExpireCtx(ctx, key, int(constant.TokenExpireTime.Seconds()))
	if err != nil {
		logc.Info(ctx, xerr.NewErrCode(xerr.KeyExpireError))
	}
}

// 删除 token

func (rds *defaultTokenModel) DeleteToken(token string) error {
	key := constant.TokenPrefixStr + token
	_, err := rds.HdelCtx(ctx, key, constant.KeyUserId, constant.KeyUserRole)
	if err != nil {
		logc.Infof(ctx, "redis HMGet key err: %v", err)
		return xerr.NewErrCode(xerr.KeyDelError)
	}
	return nil
}
