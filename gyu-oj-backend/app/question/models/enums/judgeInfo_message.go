package enums

type JudgeInfoMessage string

const (
	Accepted            JudgeInfoMessage = "Accepted"
	WrongAnswer                          = "Wrong Answer"
	CompileError                         = "Compile Error"
	MemoryLimitExceeded                  = "Memory Limit Exceeded"
	TimeLimitExceeded                    = "Time Limit Exceeded"
	PresentationError                    = "Presentation Error"
	Waiting                              = "Waiting"
	OutputLimitExceeded                  = "Output Limit Exceeded"
	DangerousOperation                   = "Dangerous Operation"
	RuntimeError                         = "Runtime Error"
	SystemError                          = "System Error"
)
