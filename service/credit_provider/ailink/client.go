package ailink

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/global"
	"go.uber.org/zap"
)

type Client struct {
	BaseURL   string
	APIKey    string
	APISecret string
	HTTP      *http.Client
}

func NewClient() *Client {
	cfg := global.GVA_CONFIG.Ailink
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseUrl), "/")
	if base == "" {
		base = DefaultBaseURL
	}
	return &Client{
		BaseURL:   base,
		APIKey:    strings.TrimSpace(cfg.APIKey),
		APISecret: strings.TrimSpace(cfg.APISecret),
		HTTP:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ListCardBins() ([]CardBin, error) {
	var out CardBinList
	if err := c.post(pathCardBins, map[string]any{}, &out); err != nil {
		return nil, err
	}
	return out.List, nil
}

func (c *Client) AccountBalances() ([]AccountBalance, error) {
	var out AccountBalanceData
	if err := c.post(pathAccountBalance, map[string]any{}, &out); err != nil {
		return nil, err
	}
	return out.Balances, nil
}

func (c *Client) CreateCardholder(req CardholderCreate) (*Cardholder, error) {
	body := map[string]any{
		"firstName": strings.TrimSpace(req.FirstName),
		"lastName":  strings.TrimSpace(req.LastName),
		"email":     strings.TrimSpace(req.Email),
		"idType":    firstNonEmpty(req.IDType, "PASSPORT"),
		"idNumber":  firstNonEmpty(strings.TrimSpace(req.IDNumber), "P1234567"),
		"state":     strings.TrimSpace(req.State),
		"city":      strings.TrimSpace(req.City),
		"address":   strings.TrimSpace(req.Address),
		"postCode":  strings.TrimSpace(req.PostCode),
		"country":   strings.TrimSpace(req.Country),
	}
	if s := strings.TrimSpace(req.RequestNo); s != "" {
		body["requestNo"] = s
	}
	if s := strings.TrimSpace(req.Mobile); s != "" {
		body["mobile"] = s
	}
	if s := strings.TrimSpace(req.Remark); s != "" {
		body["remark"] = s
	}
	var out Cardholder
	if err := c.post(pathCardholdersCreate, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ApplyCard(req ApplyCardRequest) (*ApplyCardResult, error) {
	var out ApplyCardResult
	if err := c.post(pathCardsApply, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ApplyStatus(taskID int64, batchNo, requestNo string) (*ApplyStatusResult, error) {
	body := map[string]any{}
	if taskID > 0 {
		body["taskId"] = taskID
	}
	if s := strings.TrimSpace(batchNo); s != "" {
		body["batchNo"] = s
	}
	if s := strings.TrimSpace(requestNo); s != "" {
		body["requestNo"] = s
	}
	var out ApplyStatusResult
	if err := c.post(pathCardsApplyStatus, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) QueryCards(q CardQuery) (*CardPage, error) {
	body := map[string]any{}
	if q.PageNum > 0 {
		body["pageNum"] = q.PageNum
	}
	if q.PageSize > 0 {
		body["pageSize"] = q.PageSize
	}
	putStr(body, "cardType", q.CardType)
	putStr(body, "batchNo", q.BatchNo)
	putStr(body, "cardNo", q.CardNo)
	putStr(body, "last4", q.Last4)
	putStr(body, "beginTime", q.BeginTime)
	putStr(body, "endTime", q.EndTime)
	var out CardPage
	if err := c.post(pathCardsQuery, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CardBalances(cardIDs []string) ([]CardBalance, error) {
	var out CardBalanceList
	if err := c.post(pathCardsBalances, map[string]any{"cardIds": cardIDs}, &out); err != nil {
		return nil, err
	}
	return out.List, nil
}

func (c *Client) PrivateInfo(cardID string) (*PrivateInfo, error) {
	var out PrivateInfo
	if err := c.post(pathCardsPrivate, map[string]any{"cardId": strings.TrimSpace(cardID)}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateRemark(requestNo, cardID, remark string) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathCardsRemark, map[string]any{
		"requestNo": requestNo,
		"cardId":    cardID,
		"remark":    remark,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) AdjustLimit(body map[string]any) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathCardsLimit, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) OperateStatus(requestNo, cardID, operationType string) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathCardsStatus, map[string]any{
		"requestNo":     requestNo,
		"cardId":        cardID,
		"operationType": operationType,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) OperateFunds(requestNo, cardID, operationType, amount string) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathCardsFunds, map[string]any{
		"requestNo":     requestNo,
		"cardId":        cardID,
		"operationType": operationType,
		"amount":        amount,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) QueryTransactions(q TxnQuery) (*CardTxnPage, error) {
	body := map[string]any{}
	if q.PageNum > 0 {
		body["pageNum"] = q.PageNum
	}
	if q.PageSize > 0 {
		body["pageSize"] = q.PageSize
	}
	putStr(body, "cardNo", q.CardNo)
	putStr(body, "transactionType", q.TransactionType)
	putStr(body, "transactionStatus", q.TransactionStatus)
	putStr(body, "beginTime", q.BeginTime)
	putStr(body, "endTime", q.EndTime)
	var out CardTxnPage
	if err := c.post(pathCardTransactions, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) BudgetRecharge(requestNo string, budgetID int64, amount string) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathBudgetsRecharge, map[string]any{
		"requestNo": requestNo,
		"budgetId":  budgetID,
		"amount":    amount,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) BudgetWithdraw(requestNo string, budgetID int64, amount string) (*OperateResult, error) {
	var out OperateResult
	if err := c.post(pathBudgetsWithdraw, map[string]any{
		"requestNo": requestNo,
		"budgetId":  budgetID,
		"amount":    amount,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) post(path string, body any, out any) error {
	if c.APIKey == "" || c.APISecret == "" {
		return fmt.Errorf("ailink: api-key/api-secret required")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if len(raw) == 0 || string(raw) == "null" {
		raw = []byte("{}")
	}
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce, err := randomNonce()
	if err != nil {
		return err
	}
	sig := SignRequest(http.MethodPost, path, ts, nonce, c.APISecret, raw)
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderAPIKey, c.APIKey)
	req.Header.Set(HeaderTimestamp, ts)
	req.Header.Set(HeaderNonce, nonce)
	req.Header.Set(HeaderSignature, sig)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return err
	}
	if global.GVA_LOG != nil {
		global.GVA_LOG.Info("ailink response",
			zap.String("path", path),
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(respBody)),
		)
	}
	var env apiEnvelope
	if err := json.Unmarshal(respBody, &env); err != nil {
		return fmt.Errorf("ailink %s: http %d: %s", path, resp.StatusCode, truncate(respBody, 300))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || env.Code != 200 {
		return fmt.Errorf("ailink %s: http %d code %d %s requestId=%s", path, resp.StatusCode, env.Code, env.Msg, env.RequestID)
	}
	if out == nil {
		return nil
	}
	var wrap struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(respBody, &wrap); err != nil {
		return err
	}
	if len(wrap.Data) == 0 || string(wrap.Data) == "null" {
		return nil
	}
	return json.Unmarshal(wrap.Data, out)
}

func randomNonce() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func putStr(m map[string]any, k, v string) {
	if s := strings.TrimSpace(v); s != "" {
		m[k] = s
	}
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}
	return b
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
