package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR           = 7
	AUTHREQUIRED    = 14
	KYCREQUIRED     = 15
	KYCWAITREQUIRED = 16
	SUCCESS         = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "operation_fail", c)
}
func AuthRequired(c *gin.Context) {
	Result(AUTHREQUIRED, map[string]interface{}{}, "operation_fail", c)
}
func KYCRequired(c *gin.Context) {
	Result(KYCREQUIRED, map[string]interface{}{}, "operation_fail", c)
}
func KYCWaitRequired(c *gin.Context) {
	Result(KYCWAITREQUIRED, map[string]interface{}{}, "operation_fail", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
