package gzy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
)

// SandBoxTransactionRequest POST /vcc/open/v2/sandBoxTransaction 沙箱交易模拟。
type SandBoxTransactionRequest struct {
	RequestID           string          `json:"requestId"`
	CardID              string          `json:"cardID"`
	CVV                 string          `json:"cvv"`
	ExpirationDate      string          `json:"expirationDate"` // MM/YY
	OriginTransactionID string          `json:"originTransactionId,omitempty"`
	TxnCurrency         string          `json:"txnCurrency"`
	TxnAmount           decimal.Decimal `json:"txnAmount"`
	TxnType             string          `json:"txnType"` // auth | void | refund
	MCC                 string          `json:"mcc"`
	MerchantName        string          `json:"merchantName"`
	MerchantCountry     string          `json:"merchantCountry"`
	MerchantCity        string          `json:"merchantCity"`
	MerchantPostcode    string          `json:"merchantPostcode"`
}

// SandBoxTransactionResponse Photon 沙箱模拟应答。
type SandBoxTransactionResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func validateSandBoxTransactionRequest(req *SandBoxTransactionRequest) error {
	if strings.TrimSpace(req.RequestID) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: requestId 必填")
	}
	if strings.TrimSpace(req.CardID) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: cardID 必填")
	}
	if strings.TrimSpace(req.CVV) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: cvv 必填")
	}
	if strings.TrimSpace(req.ExpirationDate) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: expirationDate 必填")
	}
	if strings.TrimSpace(req.TxnCurrency) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: txnCurrency 必填")
	}
	if !req.TxnAmount.IsPositive() {
		return fmt.Errorf("gzy sandBoxTransaction: txnAmount 须为正数")
	}
	switch strings.ToLower(strings.TrimSpace(req.TxnType)) {
	case "auth", "void", "refund":
	default:
		return fmt.Errorf("gzy sandBoxTransaction: txnType 须为 auth、void 或 refund")
	}
	if strings.TrimSpace(req.MCC) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: mcc 必填")
	}
	if strings.TrimSpace(req.MerchantName) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: merchantName 必填")
	}
	if strings.TrimSpace(req.MerchantCountry) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: merchantCountry 必填")
	}
	if strings.TrimSpace(req.MerchantCity) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: merchantCity 必填")
	}
	if strings.TrimSpace(req.MerchantPostcode) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: merchantPostcode 必填")
	}
	tt := strings.ToLower(strings.TrimSpace(req.TxnType))
	if (tt == "void" || tt == "refund") && strings.TrimSpace(req.OriginTransactionID) == "" {
		return fmt.Errorf("gzy sandBoxTransaction: void/refund 须填写 originTransactionId")
	}
	return nil
}

// EnsureAccessToken 若全局 token 缺失或即将过期则刷新（沙箱/集成测试用）。
func EnsureAccessToken() error {
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.APPID) == "" {
		return fmt.Errorf("gzy 未配置 app-id")
	}
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.AccessToken) != "" &&
		time.Now().UnixMilli() < global.GVA_CONFIG.Gzy.ExpiresAt-60_000 {
		return nil
	}
	tr, err := NewGzy().GetToken(global.GVA_CONFIG.Gzy.APPID, global.GVA_CONFIG.Gzy.APPSecret)
	if err != nil {
		return err
	}
	global.GVA_CONFIG.Gzy.AccessToken = tr.AccessToken
	global.GVA_CONFIG.Gzy.ExpiresAt = tr.ExpiresIn
	return nil
}

// ExpiryToSandBoxMMYY 将库内有效期转为沙箱要求的 MM/YY。
func ExpiryToSandBoxMMYY(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if strings.Contains(s, "/") {
		return s
	}
	if len(s) == 4 {
		for i := 0; i < 4; i++ {
			if s[i] < '0' || s[i] > '9' {
				return s
			}
		}
		return s[2:4] + "/" + s[0:2]
	}
	return s
}

// SandBoxTransaction 沙箱环境模拟发卡交易（auth / void / refund）。
func (g *Gzy) SandBoxTransaction(req SandBoxTransactionRequest) (*SandBoxTransactionResponse, error) {
	if err := validateSandBoxTransactionRequest(&req); err != nil {
		return nil, err
	}
	req.TxnType = strings.ToLower(strings.TrimSpace(req.TxnType))
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gzy sandBoxTransaction: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathSandBoxTransaction
	hreq, err := g.newRequest(http.MethodPost, reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy sandBoxTransaction: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out SandBoxTransactionResponse
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, fmt.Errorf("gzy sandBoxTransaction: decode: %w", err)
	}
	if strings.TrimSpace(out.Code) != "" && !isPhotonSuccess(out.Code) {
		return &out, gzyAPIFailure(out.Code, out.Msg)
	}
	return &out, nil
}
