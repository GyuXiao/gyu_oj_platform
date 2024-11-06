package listen

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/service"
	"gyu-oj-backend/app/judge/cmd/mq/internal/config"
	"gyu-oj-backend/app/judge/cmd/mq/internal/svc"
	"gyu-oj-backend/app/judge/cmd/rpc/judge"
	"gyu-oj-backend/common/xerr"
)

func NewListenerService(c config.Config, ctx context.Context) service.Service {
	svCtx := svc.NewServiceContext(c)
	return rabbitmq.MustNewListener(c.ListenerConf, Handler{
		Ctx:    ctx,
		SvcCtx: svCtx,
	})
}

type Handler struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func (h Handler) Consume(message string) error {
	fmt.Printf("[listener] receive: %s\n", message)
	questionSubmitId := message
	if questionSubmitId == "" {
		logc.Info(h.Ctx, "questionSubmitId 是空字符串")
		return errors.Wrapf(xerr.NewErrCode(xerr.QuestionSubmitIdIsNilError), "questionSubmitId: %s", questionSubmitId)
	}

	_, err := h.SvcCtx.JudgeRpc.DoJudge(h.Ctx, &judge.JudgeReq{
		QuestionSubmitId: questionSubmitId,
	})
	if err != nil {
		logc.Infof(h.Ctx, "调用 judge-rpc 服务时发生错误, err: %v", err)
		return err
	}

	return nil
}
