package userlogic

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"gyu-oj-backend/app/user/cmd/rpc/internal/svc"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"time"
)

type GenerateTokenReq struct {
	userId uint64
}

type GenerateTokenResp struct {
	accessToken  string
	accessExpire int64
}

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(req *GenerateTokenReq) (*GenerateTokenResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		now,
		now+accessExpire,
		req.userId,
	)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.TokenCreateFail)
	}
	return &GenerateTokenResp{
		accessToken:  accessToken,
		accessExpire: accessExpire,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds int64, userId uint64) (string, error) {
	claims := make(jwt.MapClaims)
	// 生效时间
	claims["iat"] = iat
	// 到期时间
	claims["exp"] = iat + seconds
	claims[constant.KeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *GenerateTokenLogic) ParseTokenByKey(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		switch vError, ok := err.(*jwt.ValidationError); ok {
		case vError.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, xerr.NewErrCode(xerr.TokenMalformed)
		case vError.Errors&jwt.ValidationErrorExpired != 0:
			return nil, xerr.NewErrCode(xerr.TokenExpire)
		case vError.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, xerr.NewErrCode(xerr.TokenNotValidYet)
		default:
			return nil, xerr.NewErrCode(xerr.TokenInvalid)
		}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, xerr.NewErrCode(xerr.TokenInvalid)
}
