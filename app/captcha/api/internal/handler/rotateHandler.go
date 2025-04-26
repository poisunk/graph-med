package handler

import (
	"graph-med/pkg/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"graph-med/app/captcha/api/internal/logic"
	"graph-med/app/captcha/api/internal/svc"
	"graph-med/app/captcha/api/internal/types"
)

// rotate captcha
func rotateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RotateCaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRotateLogic(r.Context(), svcCtx)
		resp, err := l.Rotate(&req)
		result.HttpResult(r, w, resp, err)
	}
}
