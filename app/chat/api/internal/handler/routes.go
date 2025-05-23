// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.2

package handler

import (
	"net/http"

	chat "graph-med/app/chat/api/internal/handler/chat"
	"graph-med/app/chat/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				// 发起对话
				Method:  http.MethodPost,
				Path:    "/chat/completion",
				Handler: chat.ChatCompletionHandler(serverCtx),
			},
			{
				// 对话反馈
				Method:  http.MethodPost,
				Path:    "/feedback",
				Handler: chat.FeedbackHandler(serverCtx),
			},
			{
				// 创建对话session
				Method:  http.MethodPost,
				Path:    "/session/create",
				Handler: chat.CreateChatSessionHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/chat/v1"),
	)
}
