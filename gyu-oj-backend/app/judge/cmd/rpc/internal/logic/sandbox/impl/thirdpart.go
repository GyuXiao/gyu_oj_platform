package impl

import (
	"gyu-oj-backend/app/judge/models/enums"
	"gyu-oj-backend/app/judge/models/types"
)

type ThirdPartSandbox struct {
}

func NewThirdPartSandbox() *ThirdPartSandbox {
	return &ThirdPartSandbox{}
}

func (sb *ThirdPartSandbox) ExecuteCode(req *types.ExecuteCodeReq) (resp *types.ExecuteCodeResp, err error) {
	return &types.ExecuteCodeResp{
		OutputList: req.InputList,
		Message:    "第三方代码沙箱--执行代码成功",
		Status:     enums.SUCCESS,
		JudgeInfo: types.JudgeInfo{
			Message: enums.Accepted,
			Time:    1000,
			Memory:  2000,
		},
	}, nil
}
