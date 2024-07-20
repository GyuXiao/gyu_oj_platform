package questionSubmit

import (
	"gyu-oj-backend/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gyu-oj-backend/app/question/cmd/api/internal/logic/questionSubmit"
	"gyu-oj-backend/app/question/cmd/api/internal/svc"
	"gyu-oj-backend/app/question/cmd/api/internal/types"
)

func CreateQuestionSubmitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateQuestionSubmitReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := questionSubmit.NewCreateQuestionSubmitLogic(r.Context(), svcCtx)
		resp, err := l.CreateQuestionSubmit(&req)
		result.HttpResult(r, w, resp, err)
	}
}
