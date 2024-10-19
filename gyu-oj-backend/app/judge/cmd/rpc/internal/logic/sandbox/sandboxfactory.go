package sandbox

import (
	"gyu-oj-backend/app/judge/cmd/rpc/internal/logic/sandbox/impl"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/judge/models/enums"
)

func SandboxFactory(ctx *svc.ServiceContext) SandboxService {
	switch ctx.Config.CodeSandbox.Type {
	case enums.Example:
		return impl.NewExampleSandbox()
	case enums.Remote:
		return impl.NewRemoteSandbox(ctx)
	case enums.ThirdPart:
		return impl.NewThirdPartSandbox()
	default:
		return impl.NewExampleSandbox()
	}
}
