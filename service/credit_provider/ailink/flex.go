package ailink

import (
	"bytes"
	"encoding/json"
)

// FlexString 同时接受 JSON 字符串和数字。卡段金额字段两种都会出现。
type FlexString string

func (f *FlexString) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*f = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	*f = FlexString(b)
	return nil
}

func (f FlexString) String() string { return string(f) }
