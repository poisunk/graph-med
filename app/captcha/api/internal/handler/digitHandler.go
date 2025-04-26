package handler

import (
	"graph-med/pkg/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"graph-med/app/captcha/api/internal/logic"
	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"
)

// digit captcha
func digitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DigitCaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDigitLogic(r.Context(), svcCtx)
		resp, err := l.Digit(&req)
		result.HttpResult(r, w, resp, err)
	}
}
