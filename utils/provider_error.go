package utils

import (
	"strings"

	"gitlab.com/ucard/service/credit_provider/gzy"
)

// ProviderUserMessage 从上游卡台（光子 / cardbin）错误中提取可展示文案。
func ProviderUserMessage(err error) string {
	if err == nil {
		return ""
	}
	if msg := gzy.UserError(err); msg != "" {
		return msg
	}
	if msg := cardbinErrorMessage(err); msg != "" {
		return msg
	}
	return ""
}

func cardbinErrorMessage(err error) string {
	s := err.Error()
	if i := strings.Index(s, "message="); i >= 0 && strings.HasPrefix(s, "cardbin API: ") {
		return strings.TrimSpace(s[i+len("message="):])
	}
	return ""
}
