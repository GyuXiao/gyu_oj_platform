package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/logic/sandbox"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/logic/strategy"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/judge/cmd/rpc/pb"
	"gyu-oj-backend/app/judge/models/enums"
	"gyu-oj-backend/app/judge/models/types"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
	pb2 "gyu-oj-backend/app/question/cmd/rpc/pb"
	"gyu-oj-backend/common/xerr"
	"time"
)

type DoJudgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDoJudgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoJudgeLogic {
	return &DoJudgeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DoJudgeLogic) DoJudge(in *pb.JudgeReq) (*pb.JudgeResp, error) {
	// 1,根据题目提交 id，获取到对应的题目、提交信息（包含代码、编程语言等）
	questionSubmitResp, err := l.svcCtx.QuestionSubmitRpc.QueryQuestionSubmitById(l.ctx, &questionsubmit.QuestionSubmitQueryByIdReq{
		Id: in.QuestionSubmitId,
	})
	if err != nil {
		logc.Infof(l.ctx, "根据 id 查询 questionSubmit 信息错误: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.QueryQuestionSubmitByIdError), "根据 id 查询 questionSubmit 错误, err: %v, questionSubmitId: %s", err, in.QuestionSubmitId)
	}
	questionResp, err := l.svcCtx.QuestionRpc.GetQuestionById(l.ctx, &question.QuestionGetByIdReq{
		Id: questionSubmitResp.QuestionSubmitVO.QuestionId,
	})
	if err != nil {
		logc.Infof(l.ctx, "根据 id 查询 question 信息错误: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.SearchQuestionByIdError), "根据 id 查询 question 错误, err: %v, questionId: %s", err, questionSubmitResp.QuestionSubmitVO.QuestionId)
	}

	// 2,如果题目提交状态不为等待中，直接返回
	if questionSubmitResp.QuestionSubmitVO.Status != enums.WAITING {
		logc.Info(l.ctx, "题目提交状态不是 waiting")
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(xerr.ServerCommonError, "题目正在在判题中，请等待一下~"), "题目提交状态: %v", questionSubmitResp.QuestionSubmitVO.Status)
	}

	// 3,更新判题（题目提交）的状态为 “判题中”，防止重复执行，也能让用户即时看到状态
	err = l.updateQuestionSubmit(in.QuestionSubmitId, nil, enums.RUNNING)
	if err != nil {
		logc.Infof(l.ctx, "更新题目提交状态为 running 错误: %v", err)
		return nil, errors.Wrap(xerr.NewErrCode(xerr.UpdateQuestionSubmitError), "更新题目提交状态错误")
	}

	// 4,调用沙箱，获取到执行结果
	sandboxImpl := sandbox.SandboxFactory(l.svcCtx)
	sandboxProxy := sandbox.SandboxProxy{RealSandbox: sandboxImpl}
	var judgeCaseInputList []string
	var judgeCaseOutputList []string
	for _, jc := range questionResp.QuestionVO.JudgeCase {
		judgeCaseInputList = append(judgeCaseInputList, jc.Input)
		judgeCaseOutputList = append(judgeCaseOutputList, jc.Output)
	}
	executeCodeResp, err := sandboxProxy.ExecuteCode(&types.ExecuteCodeReq{
		InputList: judgeCaseInputList,
		Code:      questionSubmitResp.QuestionSubmitVO.SubmitCode,
		Language:  questionSubmitResp.QuestionSubmitVO.Language,
	})
	if err != nil {
		logc.Infof(l.ctx, "调用代码沙箱错误: %v", err)
		_ = l.updateQuestionSubmit(in.QuestionSubmitId, nil, enums.FAILED)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvokeCodeSandboxError), "调用代码沙箱执行代码发生错误, err: %v", err)
	}

	// 5,根据沙箱的执行结果，设置题目的判题状态和信息
	judgeContext := &strategy.JudgeContext{
		Ctx:              l.ctx,
		ExecuteCodeResp:  executeCodeResp,
		JudgeCaseList:    judgeCaseOutputList,
		QuestionVO:       questionResp.QuestionVO,
		QuestionSubmitVO: questionSubmitResp.QuestionSubmitVO,
	}
	manager := strategy.NewJudgeManager()
	// 最终判题结果
	resp, err := manager.DoJudge(judgeContext)
	if err != nil {
		logc.Infof(l.ctx, "执行判题策略错误: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ExecuteJudgeStrategyError), "执行判题策略错误, err: %v", err)
	}

	// 6,更新数据库中的判题结果
	judgeInfo := &pb2.JudgeInfo{
		Message: resp.Message,
		Time:    resp.Time,
		Memory:  resp.Memory,
	}
	err = l.updateQuestionSubmit(in.QuestionSubmitId, judgeInfo, enums.SUCCESS)
	if err != nil {
		logc.Infof(l.ctx, "更新题目提交信息错误: %v", err)
		return nil, errors.Wrap(xerr.NewErrCode(xerr.UpdateQuestionSubmitError), "更新题目提交信息错误")
	}

	return &pb.JudgeResp{
		Id:               questionSubmitResp.QuestionSubmitVO.Id,
		Language:         questionSubmitResp.QuestionSubmitVO.Language,
		SubmitCode:       questionSubmitResp.QuestionSubmitVO.SubmitCode,
		JudgeInfoMessage: resp.Message,
		JudgeInfoTime:    resp.Time,
		JudgeInfoMemory:  resp.Memory,
		Status:           enums.SUCCESS,
		QuestionId:       questionSubmitResp.QuestionSubmitVO.QuestionId,
		UserId:           questionSubmitResp.QuestionSubmitVO.UserId,
		CreateTime:       time.Now().Unix(),
		UpdateTime:       time.Now().Unix(),
	}, nil
}

func (l *DoJudgeLogic) updateQuestionSubmit(id string, judgeInfo *pb2.JudgeInfo, status int64) error {
	req := &question.QuestionSubmitUpdateReq{
		Id:     id,
		Status: status,
	}
	if judgeInfo != nil {
		req.JudgeInfo = judgeInfo
	}
	_, err := l.svcCtx.QuestionSubmitRpc.UpdateQuestionSubmitById(l.ctx, req)

	if err != nil {
		return err
	}
	return nil
}
