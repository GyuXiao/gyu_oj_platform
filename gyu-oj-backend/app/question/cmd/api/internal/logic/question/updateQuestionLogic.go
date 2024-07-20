package question

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/question/cmd/rpc/client/question"
	"gyu-oj-backend/app/user/cmd/rpc/client/user"
	"gyu-oj-backend/common/constant"
	"gyu-oj-backend/common/xerr"
	"strings"

	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateQuestionLogic {
	return &UpdateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateQuestionLogic) UpdateQuestion(req *types.UpdateQuestionReq) (*types.UpdateQuestionResp, error) {
	// 1,非管理员不能删除题目
	token := strings.Split(req.Authorization, " ")[1]
	currentUser, _ := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if currentUser != nil && currentUser.UserRole != constant.AdminRole {
		return nil, xerr.NewErrCode(xerr.UserNotAdminError)
	}

	judgeCases := []*question.JudgeCase{}
	if len(req.JudgeCases) > 0 {
		err := copier.Copy(&judgeCases, req.JudgeCases)
		if err != nil {
			logc.Infof(l.ctx, "judgeCases 转换错误: %v\n", err)
		}
	}

	judgeConfig := &question.JudgeConfig{}
	if req.JudgeConfig != nil {
		err := copier.Copy(&judgeConfig, req.JudgeConfig)
		if err != nil {
			logc.Infof(l.ctx, "judgeConfig 转换错误: %v\n", err)
		}
	}

	// 2,调用 rpc 模块的更新题目
	resp, err := l.svcCtx.QuestionRpc.UpdateQuestion(l.ctx, &question.QuestionUpdateReq{
		Id:          req.Id,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		Answers:     req.Answers,
		JudgeCases:  judgeCases,
		JudgeConfig: judgeConfig,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateQuestionResp{IsUpdated: resp.UpdateOK}, nil
}
