package xerr

// 业务码
//前 3 位代表业务,后 3 位代表具体功能

//全局错误码

const (
	SUCCESS           uint32 = 0
	ERROR             uint32 = 1
	UnknownError      uint32 = 100000
	ServerCommonError uint32 = 100001
	ParamFormatError  uint32 = 100002
	RequestParamError uint32 = 100003
	UnauthorizedError uint32 = 100004
)

// JWT

const (
	TokenExpire            uint32 = 200001
	TokenNotValidYet       uint32 = 200002
	TokenMalformed         uint32 = 200002
	TokenInvalid           uint32 = 200003
	TokenCreateFail        uint32 = 200004
	PermissionDenied       uint32 = 200005
	NotLogin               uint32 = 200006
	LoginExpired           uint32 = 200007
	TokenParseError        uint32 = 200008
	TokenInsertError       uint32 = 200009
	TokenGetFromCacheError uint32 = 200010
)

// encryption

const (
	EncryptionError uint32 = 201001
	DecodeMd5Error  uint32 = 201002
)

// DB

const (
	RecordDuplicateError uint32 = 300001
	RecordNotFoundError  uint32 = 300002
	RecordUpdateError    uint32 = 300003
	RecordDeleteError    uint32 = 300004
	RecordCreateError    uint32 = 300005
	RecordCountError     uint32 = 300006
)

// Redis

const (
	KeyExpireError uint32 = 400001
	KeyDelError    uint32 = 400002
	KeyInsertError uint32 = 400003
)

// User

const (
	UserNotExistError          uint32 = 500001
	UserExistError             uint32 = 500002
	UserLoginError             uint32 = 500003
	UserRegisterError          uint32 = 500004
	UserPasswordError          uint32 = 500005
	UserIdNotExistError        uint32 = 500006
	UserNotLoginError          uint32 = 500007
	SearchUserError            uint32 = 500008
	CreateUserError            uint32 = 500009
	SearchUserByAccessKeyError uint32 = 500010
	AccessKeyNotExistError     uint32 = 500011
	UserNotAdminError          uint32 = 500012
)

// SDK
const (
	SDKNewClientError   uint32 = 600001
	SDKSendRequestError uint32 = 600002
)

// Question
const (
	QuestionNotExistError       uint32 = 700001
	AddQuestionError            uint32 = 700002
	UpdateQuestionError         uint32 = 700003
	DeleteQuestionError         uint32 = 700004
	SearchQuestionByIdError     uint32 = 700005
	SearchQuestionPageListError uint32 = 700006
)

// JSON
const (
	JSONMarshalError   uint32 = 800001
	JSONUnmarshalError uint32 = 800002
)

// QuestionSubmit

const (
	CreateQuestionSubmitError    uint32 = 900001
	QueryQuestionSubmitError     uint32 = 900002
	QueryQuestionSubmitByIdError uint32 = 900003
	UpdateQuestionSubmitError    uint32 = 900004
	QuestionSubmitNotExistError  uint32 = 900005
	QuestionSubmitIdIsNilError   uint32 = 900006
)

// sandbox
const (
	CompileFailError    uint32 = 110001
	RunFailError        uint32 = 110002
	RunTimeoutError     uint32 = 110003
	RunOutOfMemoryError uint32 = 110004
	SandboxError        uint32 = 110005
)

// judge
const (
	InvokeCodeSandboxError    uint32 = 120001
	ExecuteJudgeStrategyError uint32 = 120002
)

// mq
const (
	SendMessageError uint32 = 130001
)
