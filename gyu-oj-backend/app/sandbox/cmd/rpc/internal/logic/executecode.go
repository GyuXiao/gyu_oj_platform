package logic

import (
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"
	"gyu-oj-backend/app/sandbox/models"
)

// 代码沙箱的操作全流程
// 1,保存用户代码到文件中
// 2,编译代码
// 3,运行代码
// 4,整理代码的输出信息
// 5,清理编译运行时产生的额外文件

type ExecuteCodeItf interface {
	SaveCodeToFile([]byte) (string, error)
	CompileCode(string) error
	RunCode(string, []string) ([]*models.ExecResult, error)
	GetOutputResponse([]*models.ExecResult) *pb.ExecuteCodeResp
	DropFile(string) error
}
