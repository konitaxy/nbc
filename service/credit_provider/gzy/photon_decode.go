package gzy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type photonEnvelope struct {
	Code    json.RawMessage `json:"code"`
	Msg     string          `json:"msg"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func photonErrMsg(p photonEnvelope) string {
	if p.Msg != "" {
		return p.Msg
	}
	return p.Message
}

func photonCodeString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return strconv.FormatInt(int64(f), 10)
	}
	return strings.TrimSpace(string(bytes.Trim(raw, `"`)))
}

func isPhotonSuccess(code string) bool {
	c := strings.TrimSpace(strings.ToUpper(code))
	switch c {
	case "", "SUCCESS", "0000", "200", "0":
		return true
	}
	if n, err := strconv.Atoi(c); err == nil && (n == 0 || n == 200) {
		return true
	}
	return false
}

// decodePhotonJSON 解析 PhotonPay 统一应答 {code,msg,data}；若无业务 code（如 OAuth2 裸响应）则直接解析到 into。
func decodePhotonJSON(body []byte, into interface{}) error {
	var env photonEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}
	codeStr := photonCodeString(env.Code)
	hasBiz := codeStr != "" || env.Msg != "" || env.Message != ""

	if !hasBiz {
		return json.Unmarshal(body, into)
	}
	if !isPhotonSuccess(codeStr) {
		return gzyAPIFailure(codeStr, photonErrMsg(env))
	}
	if len(env.Data) > 0 && string(env.Data) != "null" {
		return json.Unmarshal(env.Data, into)
	}
	return json.Unmarshal(body, into)
}

func readBody(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}

func httpPhotonOrBodyError(status int, body []byte) error {
	if status >= 200 && status < 300 {
		return nil
	}
	var env photonEnvelope
	if json.Unmarshal(body, &env) == nil && (env.Msg != "" || env.Message != "" || len(env.Code) > 0) {
		return gzyAPIFailure(photonCodeString(env.Code), photonErrMsg(env))
	}
	const max = 512
	if len(body) > max {
		body = body[:max]
	}
	return fmt.Errorf("gzy API: HTTP %d: %s", status, string(body))
}
