package chat

import (
	"graph-med/pkg/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"graph-med/app/chat/api/internal/logic/chat"
	"graph-med/app/chat/api/internal/svc"
	"graph-med/app/chat/api/internal/types"
)

// 创建对话session
func CreateChatSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateChatSessionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewCreateChatSessionLogic(r.Context(), svcCtx)
		resp, err := l.CreateChatSession(&req)
		result.HttpResult(r, w, resp, err)
	}
}
