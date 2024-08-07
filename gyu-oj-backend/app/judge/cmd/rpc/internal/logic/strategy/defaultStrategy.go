package strategy

import (
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/judge/models/enums"
	"gyu-oj-backend/app/judge/models/types"
)

type DefaultStrategy struct {
}

func NewDefaultStrategy() *DefaultStrategy {
	return &DefaultStrategy{}
}

func (s *DefaultStrategy) DoJudge(ctx *JudgeContext) (*types.JudgeInfo, error) {
	resp := &types.JudgeInfo{
		Message: "",
		Time:    ctx.ExecuteCodeResp.JudgeInfo.Time,
		Memory:  ctx.ExecuteCodeResp.JudgeInfo.Memory,
	}

	// 0,判断代码沙箱的运行结果是否正常

	// 1,先判断沙箱执行的结果输出数量是否和预期输出数量相等
	if len(ctx.ExecuteCodeResp.OutputList) != len(ctx.JudgeCaseList) {
		logc.Infof(ctx.Ctx, "判题输出结果的个数与预期不一致, 预期是 %v, 代码沙箱的执行结果是 %v", len(ctx.JudgeCaseList), len(ctx.ExecuteCodeResp.OutputList))
		resp.Message = enums.WrongAnswer
		return resp, nil
	}

	// 2,依次判断每一项输出和预期输出是否相等
	for i, output := range ctx.ExecuteCodeResp.OutputList {
		if output != ctx.JudgeCaseList[i] {
			logc.Infof(ctx.Ctx, "第 %v 个判题样例的输出结果与预期不一致, 预期是 %s, 代码沙箱的执行结果是 %s", i+1, ctx.JudgeCaseList[i], output)
			resp.Message = enums.WrongAnswer
			return resp, nil
		}
	}

	// 3,判题题目的限制是否符合要求
	if ctx.ExecuteCodeResp.JudgeInfo.Time > ctx.QuestionVO.JudgeConfig.TimeLimit {
		logc.Infof(ctx.Ctx, "判题输出结果使用时间大于预期时间上限, 预期是 %v, 代码沙箱的时间消耗是 %v", ctx.QuestionVO.JudgeConfig.TimeLimit, ctx.ExecuteCodeResp.JudgeInfo.Time)
		resp.Message = enums.TimeLimitExceeded
		return resp, nil
	}
	if ctx.ExecuteCodeResp.JudgeInfo.Memory > ctx.QuestionVO.JudgeConfig.MemoryLimit {
		logc.Info(ctx.Ctx, "判题输出结果使用内存空间大于预期空间上限, 预期是 %v, 代码沙箱的空间消耗是 %v", ctx.QuestionVO.JudgeConfig.MemoryLimit, ctx.ExecuteCodeResp.JudgeInfo.Memory)
		resp.Message = enums.MemoryLimitExceeded
		return resp, nil
	}

	// 4,可能还有其他的异常情况，待补充

	return &types.JudgeInfo{
		Message: enums.Accepted,
		Time:    ctx.ExecuteCodeResp.JudgeInfo.Time,
		Memory:  ctx.ExecuteCodeResp.JudgeInfo.Memory,
	}, nil
}
