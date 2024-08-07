package impl

import (
	"gyu-oj-backend/app/judge/models/enums"
	"gyu-oj-backend/app/judge/models/types"
)

type ExampleSandbox struct {
}

func NewExampleSandbox() *ExampleSandbox {
	return &ExampleSandbox{}
}

func (sb *ExampleSandbox) ExecuteCode(req *types.ExecuteCodeReq) (resp *types.ExecuteCodeResp, err error) {
	// todo: 待补充真正的代码沙箱代码
	return &types.ExecuteCodeResp{
		OutputList: []string{"1", "5"},
		Message:    "示例代码沙箱--执行代码成功",
		Status:     enums.SUCCESS,
		JudgeInfo: types.JudgeInfo{
			Message: enums.Accepted,
			Time:    100,
			Memory:  100,
		},
	}, nil
}
