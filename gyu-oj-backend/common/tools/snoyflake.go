package tools

import (
	"context"
	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logc"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenId() int64 {
	id, err := flake.NextID()
	if err != nil {
		logc.Infof(context.Background(), "flake NextID failed with %s \n", err)
		panic(err)
	}

	return int64(id)
}
