package strategy

import (
	"context"
	"gyu-oj-backend/app/judge/models/types"
	"gyu-oj-backend/app/question/cmd/rpc/pb"
)

type JudgeContext struct {
	Ctx              context.Context
	ExecuteCodeResp  *types.ExecuteCodeResp
	JudgeCaseList    []string
	QuestionVO       *pb.QuestionVO
	QuestionSubmitVO *pb.QuestionSubmitVO
}
