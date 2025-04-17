package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"graph-med/internal/controller"
	"graph-med/internal/middleware"
)

type Router struct {
	authController    *controller.AuthController
	chatController    *controller.ChatController
	diseaseController *controller.DiseaseController

	enforcer *casbin.Enforcer
}

func NewRouter(
	authController *controller.AuthController,
	chatController *controller.ChatController,
	diseaseController *controller.DiseaseController,
	enforcer *casbin.Enforcer,
) *Router {
	return &Router{
		authController:    authController,
		chatController:    chatController,
		diseaseController: diseaseController,
		enforcer:          enforcer,
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {

	v1 := engine.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authController.Login)
			auth.POST("/register", r.authController.Register)
			auth.GET("/verify_captcha", r.authController.GetCaptcha)
			auth.POST("/email_captcha", r.authController.GetEmailCaptcha)
			auth.POST("/verify_captcha", r.authController.ValidateCaptcha)
			auth.POST("/logout", r.authController.Logout)
		}

		// 聊天相关路由
		chat := v1.Group("/chat")
		chat.Use(middleware.AuthMust(), middleware.ValidatePermission(r.enforcer))
		{
			chat.POST("/session/create", r.chatController.CreateSession)
			chat.POST("/completion", r.chatController.Chat)
			chat.DELETE("/session/:id", r.chatController.DeleteSession)
			chat.PUT("/session/:id", r.chatController.UpdateSession)
			chat.POST("/feedback", r.chatController.Feedback)
		}

		// 疾病库相关路由
		disease := v1.Group("/disease")
		disease.Use(middleware.AuthMust(), middleware.ValidatePermission(r.enforcer))
		{
			disease.GET("/kg/labels", r.diseaseController.GetLabels)
			disease.POST("/search", r.diseaseController.Search)
		}
	}
}
