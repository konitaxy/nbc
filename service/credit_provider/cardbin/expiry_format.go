package cardbin

import "strings"

// ExpiryYYMMToMMYY 将 cardbin 返回的有效期 YYMM（如 3006）转为 MM/YY（06/30）。
// 已含 "/" 或非法长度则原样返回。
func ExpiryYYMMToMMYY(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || strings.Contains(s, "/") {
		return s
	}
	if len(s) != 4 {
		return s
	}
	for i := 0; i < 4; i++ {
		if s[i] < '0' || s[i] > '9' {
			return s
		}
	}
	yy := s[:2]
	mm := s[2:]
	return mm + "/" + yy
}
