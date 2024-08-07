package types

type ExecuteCodeReq struct {
	InputList []string
	Code      string
	Language  string
}

type ExecuteCodeResp struct {
	OutputList []string
	Message    string
	Status     int64
	JudgeInfo  JudgeInfo
}

type JudgeInfo struct {
	Message string
	Time    int64
	Memory  int64
}
