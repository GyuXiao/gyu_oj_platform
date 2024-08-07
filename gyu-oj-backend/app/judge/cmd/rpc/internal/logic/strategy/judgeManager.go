package strategy

import (
	"gyu-oj-backend/app/judge/models/types"
)

type JudgeManager struct {
}

func NewJudgeManager() *JudgeManager {
	return &JudgeManager{}
}

func (jm *JudgeManager) DoJudge(ctx *JudgeContext) (resp *types.JudgeInfo, err error) {
	// 使用默认的判题策略
	strategy := NewDefaultStrategy()
	// 根据 ctx 选择其他判题策略，比如 语言

	return strategy.DoJudge(ctx)
}
