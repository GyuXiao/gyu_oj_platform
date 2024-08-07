package sandbox

import "gyu-oj-backend/app/judge/models/types"

type SandboxService interface {
	ExecuteCode(req *types.ExecuteCodeReq) (*types.ExecuteCodeResp, error)
}
