package sandbox

import (
	"gyu-oj-backend/app/judge/cmd/rpc/internal/logic/sandbox/impl"
	"gyu-oj-backend/app/judge/models/enums"
)

func SandboxFactory(sandboxType string) SandboxService {
	switch sandboxType {
	case enums.Example:
		return impl.NewExampleSandbox()
	case enums.Remote:
		return impl.NewRemoteSandbox()
	case enums.ThirdPart:
		return impl.NewThirdPartSandbox()
	default:
		return impl.NewExampleSandbox()
	}
}
