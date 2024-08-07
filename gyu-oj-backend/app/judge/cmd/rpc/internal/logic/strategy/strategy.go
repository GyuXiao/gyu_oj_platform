package strategy

import "gyu-oj-backend/app/judge/models/types"

type JudgeStrategy interface {
	DoJudge(ctx *JudgeContext) (*types.JudgeInfo, error)
}
