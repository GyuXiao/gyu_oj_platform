package questionSubmit

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/client/questionsubmit"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/ctxdata"
	"gyu-oj-backend/common/xerr"
	"strconv"
	"strings"

	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionSubmitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateQuestionSubmitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionSubmitLogic {
	return &CreateQuestionSubmitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionSubmitLogic) CreateQuestionSubmit(req *types.CreateQuestionSubmitReq) (*types.CreateQuestionSubmitResp, error) {
	// 1,用户必须登陆才能提交做题记录
	token := strings.Split(req.Authorization, " ")[1]
	userId := ctxdata.GetUserIdFromCtx(l.ctx)
	if userId < 0 {
		logc.Infof(l.ctx, "从 context 中获得的 userId: %v", userId)
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	currentResp, err := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{
		AuthToken: token,
		UserId:    strconv.FormatInt(userId, 10),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	if currentResp == nil {
		return nil, xerr.NewErrCode(xerr.UserNotLoginError)
	}

	// 2,调用 rpc 模块创建题目提交记录
	resp, err := l.svcCtx.QuestionSubmitRpc.DoQuestionSubmit(l.ctx, &questionsubmit.QuestionSubmitAddReq{
		Language:   req.Language,
		SubmitCode: req.SubmitCode,
		QuestionId: req.QuestionId,
		UserId:     int64(currentResp.Id),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}

	// 3,向消息队列发送消息
	err = l.SendMessage(resp.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "发送消息错误，questionSubmitId: %s", resp.Id)
	}

	return &types.CreateQuestionSubmitResp{Id: resp.Id}, nil
}

func (l *CreateQuestionSubmitLogic) SendMessage(questionSubmitId string) error {
	if questionSubmitId == "" {
		return xerr.NewErrCode(xerr.QuestionSubmitIdIsNilError)
	}

	conf := rabbitmq.RabbitSenderConf{RabbitConf: rabbitmq.RabbitConf{
		Host:     l.svcCtx.Config.RabbitMq.Host,
		Port:     l.svcCtx.Config.RabbitMq.Port,
		Username: l.svcCtx.Config.RabbitMq.Username,
		Password: l.svcCtx.Config.RabbitMq.Password,
	}, ContentType: "text/plain"}

	sender := rabbitmq.MustNewSender(conf)
	message := questionSubmitId
	err := sender.Send("oj_exchange", "oj_routingKey", []byte(message))

	if err != nil {
		logc.Infof(l.ctx, "向 RabbitMq 消息队列发消息错误: %v", err)
		return xerr.NewErrCode(xerr.SendMessageError)
	}
	return nil
}
