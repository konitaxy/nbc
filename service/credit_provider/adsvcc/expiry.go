package adsvcc

import "strings"

// NormalizeExpireDate 将渠道 expire_date 规范为 MM/YY（如 "04/29"）。
// 已是 MM/YY 则原样保留；4 位数字按 MMYY 补斜杠，不按 cardbin 的 YYMM 翻转。
func NormalizeExpireDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "-", "/")
	s = strings.ReplaceAll(s, ".", "/")
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 2 {
			mm := pad2Digits(parts[0])
			yy := strings.TrimSpace(parts[1])
			if len(yy) == 4 {
				yy = yy[2:]
			}
			if len(mm) == 2 && len(yy) == 2 {
				return mm + "/" + yy
			}
		}
		return s
	}
	if len(s) == 4 && isAllDigits(s) {
		return s[:2] + "/" + s[2:]
	}
	return s
}

func pad2Digits(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 1 && s[0] >= '0' && s[0] <= '9' {
		return "0" + s
	}
	return s
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
