package middleware

import (
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/response"
	"graph-med/internal/utils"
	"time"
)

func AuthMust() gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *utils.Claims

		code := response.SUCCESS
		token := c.Request.Header.Get("Access-Token")
		if token == "" {
			code = response.UNAUTHORIZED
		} else {
			var err error
			claims, err = utils.ParseToken(token)
			if err != nil {
				code = response.UNAUTHORIZED
			} else if time.Now().Unix() > claims.ExpiresAt.Unix() {
				code = response.UNAUTHORIZED
			}
		}

		if code != response.SUCCESS {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		if token != "" {
			claims, err := utils.ParseToken(token)
			if err != nil {
				response.Unauthorized(c)
				c.Abort()
				return
			}
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}
		c.Next()
	}
}
