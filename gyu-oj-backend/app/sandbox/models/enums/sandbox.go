package enums

const (
	SystemError    ExecuteStatus = -1
	Success        ExecuteStatus = 0
	CompileFail    ExecuteStatus = 1
	RunFail        ExecuteStatus = 2
	RunTimeout     ExecuteStatus = 3
	RunOutOfMemory ExecuteStatus = 4
)

var ExecuteStatusMsg = map[ExecuteStatus]string{
	Success:        "代码运行正常",
	CompileFail:    "代码编译失败",
	RunFail:        "代码运行失败",
	RunTimeout:     "代码运行超时",
	RunOutOfMemory: "代码运行所需内存超过限制",
	SystemError:    "系统错误，比如用户代码文件的保存等",
}

type ExecuteStatus int64

func (es ExecuteStatus) GetStatus() int64 {
	return int64(es)
}

func (es ExecuteStatus) GetMsg() string {
	return ExecuteStatusMsg[es]
}
