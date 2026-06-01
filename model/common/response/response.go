package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/utils"
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

// FailWithServiceError 上游或业务错误直接返回文案（光子 API 优先映射 code 中文说明）。
func FailWithServiceError(c *gin.Context, err error) {
	if err == nil {
		Fail(c)
		return
	}
	msg := utils.ProviderUserMessage(err)
	if msg == "" {
		msg = err.Error()
	}
	FailWithMessage(msg, c)
}

// FailWithServiceErrorUnless 命中 unless 子串时返回固定文案，否则同 FailWithServiceError。
func FailWithServiceErrorUnless(c *gin.Context, err error, unlessContains, unlessMsg string) {
	if err == nil {
		Fail(c)
		return
	}
	if unlessContains != "" && strings.Contains(err.Error(), unlessContains) {
		FailWithMessage(unlessMsg, c)
		return
	}
	FailWithServiceError(c, err)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
