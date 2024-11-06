package logic

import (
	"context"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExecuteCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteCodeLogic {
	return &ExecuteCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExecuteCodeLogic) ExecuteCode(in *pb.ExecuteCodeReq) (*pb.ExecuteCodeResp, error) {
	// 1,new 一个代码沙箱
	sandboxByGoNative := NewSandboxByGoNative()

	// 2,使用代码沙箱
	resp, err := SandboxTemplate(sandboxByGoNative, in)
	if err != nil {
		return nil, err
	}

	// 3,返回代码输出结果
	return resp, nil
}
