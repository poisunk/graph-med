package chat

import (
	"graph-med/pkg/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"graph-med/app/chat/api/internal/logic/chat"
	"graph-med/app/chat/api/internal/svc"
	"graph-med/app/chat/api/internal/types"
)

// 对话反馈
func FeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FeedbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewFeedbackLogic(r.Context(), svcCtx)
		resp, err := l.Feedback(&req)
		result.HttpResult(r, w, resp, err)
	}
}
