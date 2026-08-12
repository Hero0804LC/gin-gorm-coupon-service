package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一返回结构
type Response struct {
	Code    int         `json:"code"`    // 业务码（0=成功）
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMsg 成功返回（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Fail 失败返回
func Fail(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
		Data:    nil,
	})
}

// FailWithCode 失败返回（自定义业务码）
func FailWithCode(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
