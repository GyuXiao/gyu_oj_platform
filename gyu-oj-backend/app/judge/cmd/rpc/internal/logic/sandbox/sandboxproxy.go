package sandbox

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/judge/models/types"
)

type SandboxProxy struct {
	RealSandbox SandboxService
}

func (sp *SandboxProxy) ExecuteCode(req *types.ExecuteCodeReq) (resp *types.ExecuteCodeResp, err error) {
	ctx := context.Background()
	logc.Info(ctx, req)
	resp, err = sp.RealSandbox.ExecuteCode(req)
	if err != nil {
		logc.Errorf(ctx, "调用代码沙箱执行代码错误, err: %v", err)
	}
	logc.Info(ctx, resp)
	return
}
