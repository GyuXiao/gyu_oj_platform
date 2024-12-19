package ctxdata

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/common/constant"
)

func GetUserIdFromCtx(ctx context.Context) int64 {
	jsonUid, ok := ctx.Value(constant.KeyJwtUserId).(json.Number)
	if !ok {
		logc.Info(ctx, "GetUserIdFromCtx err: userId 不存在")
		return -1
	}

	int64Uid, err := jsonUid.Int64()
	if err != nil {
		logc.Infof(ctx, "jsonUid.Int64() err: %+v", err)
		return -1
	}

	return int64Uid
}
