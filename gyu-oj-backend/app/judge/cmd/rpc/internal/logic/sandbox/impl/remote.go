package impl

import (
	"context"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/judge/models/types"
	"gyu-oj-backend/app/sandbox/cmd/rpc/codesandbox"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"
)

type RemoteSandbox struct {
	sandboxRpcServer codesandbox.CodeSandbox
}

func NewRemoteSandbox(ctx *svc.ServiceContext) *RemoteSandbox {
	return &RemoteSandbox{sandboxRpcServer: ctx.SandboxRpc}
}

func (sb *RemoteSandbox) ExecuteCode(req *types.ExecuteCodeReq) (resp *types.ExecuteCodeResp, err error) {
	execResp, err := sb.sandboxRpcServer.ExecuteCode(context.Background(), &pb.ExecuteCodeReq{
		InputList: req.InputList,
		Code:      req.Code,
		Language:  req.Language,
	})
	if err != nil {
		return nil, err
	}
	return &types.ExecuteCodeResp{
		OutputList: execResp.OutputList,
		Message:    execResp.Message,
		Status:     execResp.Status,
		JudgeInfo: types.JudgeInfo{
			Message: execResp.ExecuteResultMessage,
			Time:    execResp.ExecuteResultTime,
			Memory:  execResp.ExecuteResultMemory,
		},
	}, nil
}
