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

type CreateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionLogic) CreateQuestion(req *types.CreateQuestionReq) (*types.CreateQuestionResp, error) {
	// 1,非管理员不能创建题目
	token := strings.Split(req.Authorization, " ")[1]
	currentUser, _ := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if currentUser != nil && currentUser.UserRole != constant.AdminRole {
		return nil, xerr.NewErrCode(xerr.UserNotAdminError)
	}

	// 2,创建问题前进行业务参数校验
	if req.Title == "" || req.Content == "" {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "题目的标题或内容不能为空")
	}
	if req.Answer == "" {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "管理员创建的题目不能没有答案")
	}
	if len(req.Tags) == 0 {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "题目需要至少一个标签")
	}
	if len(req.JudgeCase) == 0 {
		return nil, xerr.NewErrCodeMsg(xerr.ParamFormatError, "题目需要至少一个测试用例")
	}

	judgeCases := []*question.JudgeCase{}
	err := copier.Copy(&judgeCases, req.JudgeCase)
	if err != nil {
		logc.Infof(l.ctx, "judgeCases 转换错误: %v\n", err)
	}

	judgeConfig := &question.JudgeConfig{}
	err = copier.Copy(&judgeConfig, req.JudgeConfig)
	if err != nil {
		logc.Infof(l.ctx, "judgeConfig 转换错误: %v\n", err)
	}

	// 调用 rpc 模块创建题目
	resp, err := l.svcCtx.QuestionRpc.AddQuestion(l.ctx, &question.QuestionAddReq{
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		Answer:      req.Answer,
		JudgeCases:  judgeCases,
		JudgeConfig: judgeConfig,
		UserId:      int64(currentUser.Id),
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateQuestionResp{Id: resp.Id}, nil
}
