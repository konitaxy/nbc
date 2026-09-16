package adsvcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ucard/global"
)

const DefaultBaseURL = "http://18.162.211.187:89/v1"

// Client Adsvcc OpenAPI 客户端（RSA 加签 + token）。
type Client struct {
	BaseURL    string
	AppID      string
	AppSecret  string
	PrivateKey string
	Token      string
	HTTP       *http.Client
}

func NewClient() *Client {
	cfg := global.GVA_CONFIG.Adsvcc
	base := strings.TrimRight(strings.TrimSpace(cfg.BaseUrl), "/")
	if base == "" {
		base = DefaultBaseURL
	}
	pem, _ := loadPrivateKeyPEM(cfg.PrivateKey, cfg.PrivateKeyPath)
	hc := &http.Client{Timeout: 45 * time.Second}
	if global.GVA_CONFIG.System.Env == "dev" {
		if proxyURL, err := url.Parse("http://127.0.0.1:7890"); err == nil {
			hc.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}
	return &Client{
		BaseURL:    base,
		AppID:      strings.TrimSpace(cfg.APPID),
		AppSecret:  strings.TrimSpace(cfg.APPSecret),
		PrivateKey: pem,
		Token:      strings.TrimSpace(cfg.AccessToken),
		HTTP:       hc,
	}
}

func (c *Client) syncTokenFromConfig() {
	c.Token = strings.TrimSpace(global.GVA_CONFIG.Adsvcc.AccessToken)
}

func (c *Client) GetAccessToken() (*TokenData, error) {
	if c.AppID == "" || c.AppSecret == "" {
		return nil, fmt.Errorf("adsvcc GetAccessToken: app-id/app-secret required")
	}
	form := map[string]string{
		"AppId":     c.AppID,
		"AppSecret": c.AppSecret,
	}
	var out TokenData
	if err := c.doMultipart(http.MethodPost, pathAccessToken, form, false, &out); err != nil {
		return nil, err
	}
	if strings.TrimSpace(out.Token) == "" {
		return nil, fmt.Errorf("adsvcc GetAccessToken: empty token")
	}
	return &out, nil
}

func (c *Client) RefreshAccessToken() (*TokenData, error) {
	c.syncTokenFromConfig()
	if c.Token == "" {
		return c.GetAccessToken()
	}
	var out TokenData
	if err := c.doForm(http.MethodPost, pathRefreshToken, map[string]string{}, true, &out); err != nil {
		return nil, err
	}
	if strings.TrimSpace(out.Token) == "" {
		return nil, fmt.Errorf("adsvcc RefreshAccessToken: empty token")
	}
	return &out, nil
}

func (c *Client) DemandCard(req DemandRequest) (*DemandData, error) {
	form := map[string]string{
		"product_id":         strings.TrimSpace(req.ProductID),
		"amount":             strings.TrimSpace(req.Amount),
		"num":                firstNonEmpty(strings.TrimSpace(req.Num), "1"),
		"user_cardholder_id": strconv.FormatInt(req.UserCardholderID, 10),
	}
	pt := req.ProductType
	if pt == 0 {
		pt = ProductTypeDebit
	}
	form["product_type"] = strconv.Itoa(pt)
	if mx := strings.TrimSpace(req.MatrixAccount); mx != "" {
		form["matrix_account"] = mx
	}
	var out DemandData
	if err := c.doForm(http.MethodPost, pathDemand, form, true, &out); err != nil {
		return nil, err
	}
	if strings.TrimSpace(out.BatchID) == "" {
		return nil, fmt.Errorf("adsvcc DemandCard: empty batch_id")
	}
	return &out, nil
}

func (c *Client) DemandBatchStatus(batchID string) (*DemandStatusData, error) {
	q := url.Values{}
	q.Set("batch_id", strings.TrimSpace(batchID))
	var out DemandStatusData
	if err := c.doQuery(http.MethodGet, pathDemandStatus, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCardInfo(cardID string) (*CardInfoData, error) {
	q := url.Values{}
	q.Set("card_id", strings.TrimSpace(cardID))
	var out CardInfoData
	if err := c.doQuery(http.MethodGet, pathCardInfo, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetCardCVV(cardID, password, mfaCode string) (*CardCVVData, error) {
	form := map[string]string{
		"card_id":  strings.TrimSpace(cardID),
		"password": password,
		"code":     mfaCode,
	}
	var out CardCVVData
	if err := c.doMultipart(http.MethodPost, pathCardCVV, form, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateCard(req CardUpdateRequest) error {
	form := map[string]string{
		"action":  strings.TrimSpace(req.Action),
		"card_id": strings.TrimSpace(req.CardID),
	}
	if a := strings.TrimSpace(req.Amount); a != "" {
		form["amount"] = a
	}
	if t := strings.TrimSpace(req.Tag); t != "" {
		form["tag"] = t
	}
	var ignore ignoreData
	return c.doForm(http.MethodPost, pathCardUpdate, form, true, &ignore)
}

func (c *Client) ListCards(req CardListRequest) (*CardListData, error) {
	q := url.Values{}
	if req.Page > 0 {
		q.Set("page", strconv.Itoa(req.Page))
	}
	if req.Limit > 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.CardID != "" {
		q.Set("card_id", req.CardID)
	}
	if req.Status != "" {
		q.Set("status", req.Status)
	}
	if req.Tag != "" {
		q.Set("tag", req.Tag)
	}
	var out CardListData
	if err := c.doQuery(http.MethodGet, pathCardList, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListTransactions(req TxnListRequest) (*TxnListData, error) {
	q := url.Values{}
	if req.CardID != "" {
		q.Set("card_id", req.CardID)
	}
	if req.Page > 0 {
		q.Set("page", strconv.Itoa(req.Page))
	}
	if req.Limit > 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.StartDate != "" {
		q.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		q.Set("end_date", req.EndDate)
	}
	if req.Status != "" {
		q.Set("status", req.Status)
	}
	var out TxnListData
	if err := c.doQuery(http.MethodGet, pathCardTxn, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) AddCardUse(req CardUseRequest) (*CardUseData, error) {
	form := cardUseForm(req, false)
	var raw json.RawMessage
	if err := c.doMultipart(http.MethodPost, pathCardUseAdd, form, true, &raw); err != nil {
		return nil, err
	}
	id := ParseCardUseID(raw)
	if id == "" && strings.TrimSpace(req.Email) != "" {
		if listed, err := c.ListCardUse(CardUseListRequest{Page: 1, Limit: 50}); err == nil && listed != nil {
			want := strings.ToLower(strings.TrimSpace(req.Email))
			for _, it := range listed.List {
				if strings.ToLower(strings.TrimSpace(it.Email)) == want {
					id = strings.TrimSpace(it.ID.String())
					break
				}
			}
		}
	}
	if id == "" {
		return nil, fmt.Errorf("adsvcc AddCardUse: empty holder id")
	}
	return &CardUseData{ID: id}, nil
}

func (c *Client) EditCardUse(req CardUseRequest) error {
	if req.ID <= 0 {
		return fmt.Errorf("adsvcc EditCardUse: id required")
	}
	form := cardUseForm(req, true)
	var ignore ignoreData
	return c.doMultipart(http.MethodPost, pathCardUseEdit, form, true, &ignore)
}

func (c *Client) RechargeShareWallet(req ShareWalletRequest) (*ShareRechargeData, error) {
	form := shareWalletForm(req)
	var out ShareRechargeData
	if err := c.doMultipart(http.MethodPost, pathShareRecharge, form, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ReduceShareWallet(req ShareWalletRequest) error {
	form := shareWalletForm(req)
	var ignore ignoreData
	return c.doMultipart(http.MethodPost, pathShareReduce, form, true, &ignore)
}

func (c *Client) GetShareWalletBalance(req ShareWalletRequest) (*ShareBalanceData, error) {
	q := url.Values{}
	if req.IsAuto != 0 {
		q.Set("is_auto", strconv.Itoa(req.IsAuto))
	} else {
		q.Set("is_auto", "0")
	}
	if req.ProviderID > 0 {
		q.Set("provider_id", strconv.FormatInt(req.ProviderID, 10))
	}
	var out ShareBalanceData
	if err := c.doQuery(http.MethodGet, pathShareBalance, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func shareWalletForm(req ShareWalletRequest) map[string]string {
	form := map[string]string{
		"amount": strings.TrimSpace(req.Amount),
	}
	if req.ProviderID > 0 {
		form["provider_id"] = strconv.FormatInt(req.ProviderID, 10)
	}
	return form
}

func (c *Client) ListCardUse(req CardUseListRequest) (*CardUseListData, error) {
	q := url.Values{}
	if req.Page > 0 {
		q.Set("page", strconv.Itoa(req.Page))
	}
	if req.Limit > 0 {
		q.Set("limit", strconv.Itoa(req.Limit))
	}
	var out CardUseListData
	if err := c.doQuery(http.MethodGet, pathCardUseList, q, true, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func cardUseForm(req CardUseRequest, withID bool) map[string]string {
	pt := req.Type
	if pt == 0 {
		pt = ProductTypeDebit
	}
	form := map[string]string{
		"first_name":   strings.TrimSpace(req.FirstName),
		"last_name":    strings.TrimSpace(req.LastName),
		"email":        strings.TrimSpace(req.Email),
		"phone":        strings.TrimSpace(req.Phone),
		"phone_prefix": normalizePhonePrefix(req.PhonePrefix),
		"birthday":     strings.TrimSpace(req.Birthday),
		"country_code": strings.ToUpper(strings.TrimSpace(req.CountryCode)),
		"type":         strconv.Itoa(pt),
	}
	if s := strings.TrimSpace(req.City); s != "" {
		form["city"] = s
	}
	if s := strings.TrimSpace(req.Province); s != "" {
		form["province"] = s
	}
	if s := strings.TrimSpace(req.Address); s != "" {
		form["address"] = s
	}
	if s := strings.TrimSpace(req.Zip); s != "" {
		form["zip"] = s
	}
	if withID && req.ID > 0 {
		form["id"] = strconv.FormatInt(req.ID, 10)
	}
	return form
}

func normalizePhonePrefix(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "+")
	return strings.TrimSpace(s)
}

// ParseCardUseID 从添加用卡人 data 中解析 id（文档响应 schema 为空，兼容对象/数字/字符串）。
func ParseCardUseID(data json.RawMessage) string {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return ""
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err == nil && len(obj) > 0 {
		for _, k := range []string{"id", "user_cardholder_id", "user_card_use_id"} {
			if v, ok := obj[k]; ok {
				if s := jsonScalarString(v); s != "" {
					return s
				}
			}
		}
	}
	return jsonScalarString(data)
}

func jsonScalarString(v json.RawMessage) string {
	v = bytes.TrimSpace(v)
	if len(v) == 0 || string(v) == "null" {
		return ""
	}
	var n json.Number
	if err := json.Unmarshal(v, &n); err == nil && strings.TrimSpace(n.String()) != "" {
		return strings.TrimSpace(n.String())
	}
	var s string
	if err := json.Unmarshal(v, &s); err == nil {
		return strings.TrimSpace(s)
	}
	return ""
}

type ignoreData struct{}

func (ignoreData) UnmarshalJSON([]byte) error { return nil }

func (c *Client) doForm(method, path string, form map[string]string, withToken bool, out any) error {
	c.syncTokenFromConfig()
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sign, err := SignParams(c.PrivateKey, form, ts)
	if err != nil {
		return err
	}
	body := url.Values{}
	for k, v := range form {
		body.Set(k, v)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setCommonHeaders(req, ts, sign, withToken)
	return c.decode(req, out)
}

func (c *Client) doMultipart(method, path string, form map[string]string, withToken bool, out any) error {
	c.syncTokenFromConfig()
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sign, err := SignParams(c.PrivateKey, form, ts)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range form {
		_ = w.WriteField(k, v)
	}
	_ = w.Close()
	req, err := http.NewRequest(method, c.BaseURL+path, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	c.setCommonHeaders(req, ts, sign, withToken)
	return c.decode(req, out)
}

func (c *Client) doQuery(method, path string, q url.Values, withToken bool, out any) error {
	c.syncTokenFromConfig()
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sign, err := SignParams(c.PrivateKey, nil, ts)
	if err != nil {
		return err
	}
	u := c.BaseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return err
	}
	c.setCommonHeaders(req, ts, sign, withToken)
	return c.decode(req, out)
}

func (c *Client) setCommonHeaders(req *http.Request, timestamp, sign string, withToken bool) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("sign", sign)
	if withToken && c.Token != "" {
		req.Header.Set("token", c.Token)
	}
}

func (c *Client) decode(req *http.Request, out any) error {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("adsvcc http: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("adsvcc read body: %w", err)
	}
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("adsvcc decode envelope: %w; body=%s", err, truncate(raw, 256))
	}
	msg := strings.TrimSpace(firstNonEmpty(env.Msg, env.Message))
	if env.Code != 0 {
		return fmt.Errorf("adsvcc API: code=%d, message=%s", env.Code, msg)
	}
	if out == nil {
		return nil
	}
	data := envelopePayload(env)
	if len(data) == 0 || string(data) == "null" || string(data) == "[]" {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("adsvcc decode data: %w; data=%s", err, truncate(data, 256))
	}
	return nil
}

func envelopePayload(env Envelope) json.RawMessage {
	data := bytes.TrimSpace(env.Data)
	if len(data) > 0 && string(data) != "null" {
		return data
	}
	return bytes.TrimSpace(env.Datas)
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}

// ExpiresAtMillis 将 expiresIn（秒级绝对 unix 时间戳）转为毫秒。
func ExpiresAtMillis(expiresIn json.Number) int64 {
	s := strings.TrimSpace(expiresIn.String())
	if s == "" {
		return time.Now().Add(2 * time.Hour).UnixMilli()
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now().Add(2 * time.Hour).UnixMilli()
	}
	if v < 1e12 {
		return v * 1000
	}
	return v
}
