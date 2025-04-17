package middleware

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/redis"
	"graph-med/internal/base/response"
	"graph-med/internal/model"
)

func ValidatePermission(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Access-Token")

		user, err := redis.Get[model.ChatUser]("user-" + token)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		valid, err := enforcer.Enforce(user.UserID, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			response.Error(c, errors.New("failed to validate permission: "+err.Error()))
			c.Abort()
			return
		}

		if !valid {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
