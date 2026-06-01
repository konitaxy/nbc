package gzy

import (
	"regexp"
	"strings"
)

var gzyAPICodeMsgRe = regexp.MustCompile(`^gzy API: code=([^,]+), message=(.+)$`)

// UserError 从光子（GZY）相关 error 中提取可返回前端的文案。
func UserError(err error) string {
	if err == nil {
		return ""
	}
	s := err.Error()

	if m := gzyAPICodeMsgRe.FindStringSubmatch(s); len(m) == 3 {
		return strings.TrimSpace(m[2])
	}
	if strings.HasPrefix(s, "gzy API: message=") {
		return strings.TrimSpace(s[len("gzy API: message="):])
	}
	if strings.HasPrefix(s, "gzy API: HTTP ") {
		if i := strings.LastIndex(s, ": "); i >= 0 {
			return strings.TrimSpace(s[i+2:])
		}
	}
	return gzyPrefixedMessage(s)
}

func gzyPrefixedMessage(s string) string {
	if !strings.HasPrefix(s, "gzy ") {
		return ""
	}
	if idx := strings.Index(s, ": "); idx > 0 {
		return strings.TrimSpace(s[idx+2:])
	}
	return strings.TrimSpace(strings.TrimPrefix(s, "gzy "))
}
