package v1

import (
	"net/http"

	"github.com/ewangplay/serval/io"
	"github.com/gin-gonic/gin"
)

const (
	ERROR   = -1
	SUCCESS = 0
)

func SuccResult(code int, data any, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, io.Response{
		code,
		data,
		msg,
	})
}

func FailResult(status int, code int, data any, msg string, c *gin.Context) {
	// 开始时间
	c.AbortWithStatusJSON(status, io.Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	SuccResult(SUCCESS, map[string]any{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	SuccResult(SUCCESS, map[string]any{}, message, c)
}

func OkWithData(data any, c *gin.Context) {
	SuccResult(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data any, message string, c *gin.Context) {
	SuccResult(SUCCESS, data, message, c)
}

func Fail(status int, c *gin.Context) {
	FailResult(status, ERROR, map[string]any{}, "操作失败", c)
}

func FailWithMessage(status int, message string, c *gin.Context) {
	FailResult(status, ERROR, map[string]any{}, message, c)
}

func FailWithDetailed(status int, data any, message string, c *gin.Context) {
	FailResult(status, ERROR, data, message, c)
}
