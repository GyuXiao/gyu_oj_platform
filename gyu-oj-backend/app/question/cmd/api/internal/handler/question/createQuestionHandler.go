package question

import (
	"gyu-oj-backend/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gyu-oj-backend/app/question/cmd/api/internal/logic/question"
	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"
)

func CreateQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateQuestionReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := question.NewCreateQuestionLogic(r.Context(), svcCtx)
		resp, err := l.CreateQuestion(&req)
		result.HttpResult(r, w, resp, err)
	}
}
