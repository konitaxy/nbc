package logredact

import (
	"encoding/json"
	"regexp"
	"strings"
)

const Redacted = "****"

// 兼容旧名
const CVV = Redacted

var sensitiveJSONFieldRe = regexp.MustCompile(
	`(?i)("(?:cvv|expirationDate|expiration_date|expiry|expirey)"\s*:\s*)"[^"]*"`,
)

// CardSensitiveInJSON 将 JSON 中的卡敏感字段（cvv、有效期等）替换为 ****。
func CardSensitiveInJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return sensitiveJSONFieldRe.ReplaceAllString(raw, `${1}"`+Redacted+`"`)
	}
	redactSensitiveWalk(v)
	b, err := json.Marshal(v)
	if err != nil {
		return sensitiveJSONFieldRe.ReplaceAllString(raw, `${1}"`+Redacted+`"`)
	}
	return string(b)
}

// CVVInJSON 兼容旧调用；同时脱敏 expirationDate / expiry 等字段。
func CVVInJSON(raw string) string {
	return CardSensitiveInJSON(raw)
}

func redactSensitiveWalk(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, val := range x {
			if isSensitiveCardJSONKey(k) {
				x[k] = Redacted
			} else {
				redactSensitiveWalk(val)
			}
		}
	case []interface{}:
		for i := range x {
			redactSensitiveWalk(x[i])
		}
	}
}

func isSensitiveCardJSONKey(k string) bool {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case "cvv", "expirationdate", "expiration_date", "expiry", "expirey":
		return true
	default:
		return false
	}
}
