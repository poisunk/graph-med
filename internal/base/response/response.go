package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	SUCCESS      = 200
	UNAUTHORIZED = 401
	FORBIDDEN    = 403
	ERROR        = 500
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: SUCCESS,
		Data: data,
		Msg:  "",
	})
}

func Error(c *gin.Context, err error) {
	c.JSON(200, Response{
		Code: ERROR,
		Data: nil,
		Msg:  err.Error(),
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(200, Response{
		Code: UNAUTHORIZED,
		Data: nil,
		Msg:  "Unauthorized",
	})
}

func Forbidden(c *gin.Context) {
	c.JSON(200, Response{
		Code: FORBIDDEN,
		Data: nil,
		Msg:  "Forbidden",
	})
}

func EventStream(c *gin.Context, data string) error {
	event := fmt.Sprintf("data: %s\n\n", data)
	_, err := c.Writer.Write([]byte(event))
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return err
	}
	c.Writer.Flush()
	return nil
}
