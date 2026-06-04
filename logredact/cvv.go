package logredact

import (
	"encoding/json"
	"regexp"
	"strings"
)

const CVV = "****"

var cvvJSONFieldRe = regexp.MustCompile(`(?i)("cvv"\s*:\s*)"[^"]*"`)

// CVVInJSON 将 JSON 文本中的 cvv 字段值替换为 ****。
func CVVInJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return cvvJSONFieldRe.ReplaceAllString(raw, `${1}"`+CVV+`"`)
	}
	redactCVVWalk(v)
	b, err := json.Marshal(v)
	if err != nil {
		return cvvJSONFieldRe.ReplaceAllString(raw, `${1}"`+CVV+`"`)
	}
	return string(b)
}

func redactCVVWalk(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, val := range x {
			if strings.EqualFold(strings.TrimSpace(k), "cvv") {
				x[k] = CVV
			} else {
				redactCVVWalk(val)
			}
		}
	case []interface{}:
		for i := range x {
			redactCVVWalk(x[i])
		}
	}
}
