package chat

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"graph-med/app/chat/api/internal/logic/chat"
	"graph-med/app/chat/api/internal/svc"
	"graph-med/app/chat/api/internal/types"
)

// 发起对话
func ChatCompletionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatCompletionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		l := chat.NewChatCompletionLogic(r.Context(), svcCtx)
		resultChan, err := l.ChatCompletion(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for resp := range resultChan {
			_, _ = fmt.Fprintf(w, "data: %s\n\n", resp)
			flusher.Flush()
		}

		_, _ = fmt.Fprintf(w, "data: [DONE]\n\n")
		flusher.Flush()
	}
}
