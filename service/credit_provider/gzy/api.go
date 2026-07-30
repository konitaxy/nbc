package gzy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/logredact"
	"go.uber.org/zap"
)

// DEFAULT_BASE_URL 为网关根域名；业务路径见 paths.go。取令牌默认 GET /oauth2/token/accessToken。
// 官方文档: https://api-doc.photonpay.com/
const DEFAULT_BASE_URL = "https://x-api.photonpay.com"

func gzyAPIFailure(code, message string) error {
	code = strings.TrimSpace(code)
	userMsg := photonUserMessage(code, message)
	if code != "" {
		return fmt.Errorf("gzy API: code=%s, message=%s", strings.ToUpper(code), userMsg)
	}
	return fmt.Errorf("gzy API: message=%s", userMsg)
}

type Gzy struct {
	BaseURL     string
	AccessToken string
	client      *http.Client
}

func NewGzy() *Gzy {
	client := &http.Client{}

	// dev 环境使用本地代理
	if global.GVA_CONFIG.System.Env == "dev" {
		proxyURL, _ := url.Parse("http://127.0.0.1:7890")
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 30 * time.Second,
		}
	}

	return &Gzy{
		BaseURL:     global.GVA_CONFIG.Gzy.BaseUrl,
		AccessToken: global.GVA_CONFIG.Gzy.AccessToken,
		client:      client,
	}
}

func (g *Gzy) GetToken(appID, appSecret string) (*TokenResponse, error) {
	appID = strings.TrimSpace(appID)
	appSecret = strings.TrimSpace(appSecret)
	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("gzy GetToken: appId 与 appSecret 必填")
	}
	tokenPath := strings.TrimSpace(global.GVA_CONFIG.Gzy.TokenPath)
	if tokenPath == "" {
		tokenPath = pathOAuthAccessToken
	}
	if !strings.HasPrefix(tokenPath, "/") {
		tokenPath = "/" + tokenPath
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + tokenPath
	norm := strings.TrimSuffix(strings.ToLower(tokenPath), "/")
	if norm == strings.TrimSuffix(strings.ToLower(pathOAuthTokenFormLegacy), "/") {
		return g.getTokenOAuthFormLegacy(appID, appSecret, reqURL)
	}
	return g.getTokenAccessTokenPOST(appID, appSecret, reqURL)
}

func (g *Gzy) getTokenAccessTokenPOST(appID, appSecret, reqURL string) (*TokenResponse, error) {
	basic := base64.StdEncoding.EncodeToString([]byte(appID + "/" + appSecret))
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gzy GetToken: %w", err)
	}
	req.Header.Set("Authorization", "basic "+basic)
	req.Header.Set("Request-Id", fmt.Sprintf("%d", time.Now().UnixNano()))

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gzy GetToken: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env photonAccessTokenEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy GetToken: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy GetToken: 响应缺少 data")
	}
	var d photonAccessTokenData
	if err := json.Unmarshal(env.Data, &d); err != nil {
		return nil, fmt.Errorf("gzy GetToken: decode data: %w", err)
	}
	tr := &TokenResponse{
		AccessToken:           strings.TrimSpace(d.Token),
		RefreshToken:          strings.TrimSpace(d.RefreshToken),
		ExpiresIn:             d.ExpiresIn,
		RefreshTokenExpiresIn: d.RefreshExpiresIn,
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("gzy GetToken: data.token 为空")
	}
	normalizeTokenResponseExpiry(tr)
	return tr, nil
}

func (g *Gzy) getTokenOAuthFormLegacy(appID, appSecret, reqURL string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", appID)
	form.Set("client_secret", appSecret)
	req, err := g.newPublicRequest(http.MethodPost, reqURL, []byte(form.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return nil, fmt.Errorf("gzy GetToken (legacy): %w", err)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gzy GetToken (legacy): %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var tr TokenResponse
	if err := decodePhotonJSON(body, &tr); err != nil {
		return nil, err
	}
	normalizeTokenResponseExpiry(&tr)
	return &tr, nil
}

func normalizeTokenResponseExpiry(tr *TokenResponse) {
	// 与 cardbin 侧约定一致：ExpiresIn / RefreshTokenExpiresIn 存绝对时间毫秒；接口常返回相对秒
	if tr.ExpiresIn > 0 && tr.ExpiresIn < 1e12 {
		tr.ExpiresIn = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second).UnixMilli()
	}
	if tr.RefreshTokenExpiresIn > 0 && tr.RefreshTokenExpiresIn < 1e12 {
		tr.RefreshTokenExpiresIn = time.Now().Add(time.Duration(tr.RefreshTokenExpiresIn) * time.Second).UnixMilli()
	}
}

// photonAccessTokenEnvelope GET /oauth2/token/accessToken 应答外层。
type photonAccessTokenEnvelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type photonAccessTokenData struct {
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	Token            string `json:"token"`
}

func (g *Gzy) newPublicRequest(method, urlStr string, body []byte, contentType string) (req *http.Request, err error) {
	if body == nil {
		req, err = http.NewRequest(method, urlStr, nil)
	} else {
		req, err = http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Request-Id", fmt.Sprintf("%d", time.Now().UnixNano()))
	return req, nil
}

func (g *Gzy) newRequest(method, reqURL string, body []byte) (req *http.Request, err error) {
	if strings.TrimSpace(g.AccessToken) == "" {
		return nil, fmt.Errorf("gzy: X-PD-TOKEN 为空，请先获取访问令牌")
	}
	signPayload := body
	if signPayload == nil {
		signPayload = []byte{}
	}
	sign, err := buildXPDSign(signPayload)
	if err != nil {
		return nil, err
	}

	if body == nil {
		req, err = http.NewRequest(method, reqURL, nil)
	} else {
		req, err = http.NewRequest(method, reqURL, bytes.NewBuffer(body))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-PD-TOKEN", g.AccessToken)
	if strings.TrimSpace(sign) != "" {
		req.Header.Set("X-PD-SIGN", sign)
	}
	req.Header.Set("Request-Id", fmt.Sprintf("%d", time.Now().UnixNano()))
	return req, nil
}

func httpHeadersMap(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, vs := range h {
		out[k] = strings.Join(vs, ",")
	}
	return out
}

func logGzyOpenCardRequest(hreq *http.Request, body []byte) {
	global.GVA_LOG.Info("gzy openCard HTTP request",
		zap.String("method", hreq.Method),
		zap.String("url", hreq.URL.String()),
		zap.Any("headers", httpHeadersMap(hreq.Header)),
		zap.String("body", string(body)),
	)
}

func logGzyOpenCardResponse(statusCode int, body []byte, err error) {
	fields := []zap.Field{
		zap.Int("statusCode", statusCode),
		zap.String("body", logredact.CVVInJSON(string(body))),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	global.GVA_LOG.Info("gzy openCard HTTP response", fields...)
}

func logGzyOpenCardParsed(partnerOrderID string, env *openCardV4Envelope, out *CreateCardResponse) {
	if env == nil {
		return
	}
	fields := []zap.Field{
		zap.String("partnerOrderId", partnerOrderID),
		zap.String("code", env.Code),
		zap.String("msg", env.Msg),
	}
	if env.Data != nil {
		fields = append(fields,
			zap.String("requestId", strings.TrimSpace(env.Data.RequestID)),
			zap.String("dataStatus", strings.TrimSpace(env.Data.Status)),
		)
		if d := env.Data.CardDetail; d != nil {
			fields = append(fields,
				zap.String("cardId", strings.TrimSpace(d.CardID)),
				zap.String("cardNo", strings.TrimSpace(d.CardNo)),
				zap.String("cardStatus", strings.TrimSpace(d.CardStatus)),
			)
		}
	}
	if out != nil {
		fields = append(fields, zap.String("cardId", out.CardID))
	}
	global.GVA_LOG.Info("gzy openCard parsed result", fields...)
}

// walletAccountSingleV4Envelope GET /wallet/openApi/v4/account/single 应答外层。
type walletAccountSingleV4Envelope struct {
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
	Method string          `json:"method,omitempty"`
	Path   string          `json:"path,omitempty"`
}

// GetWalletAccountSingle GET /wallet/openApi/v4/account/single 查询光子易钱包单账户实时余额。
func (g *Gzy) GetWalletAccountSingle(req WalletAccountSingleRequest) (*BalanceResponse, error) {
	q, err := buildWalletAccountSingleQuery(req)
	if err != nil {
		return nil, fmt.Errorf("gzy GetWalletAccountSingle: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathAccountSingle + "?" + q.Encode()
	hreq, err := g.newRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy GetWalletAccountSingle: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env walletAccountSingleV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy GetWalletAccountSingle: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy GetWalletAccountSingle: 响应缺少 data")
	}
	var out BalanceResponse
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, fmt.Errorf("gzy GetWalletAccountSingle: decode data: %w", err)
	}
	return &out, nil
}

// GetBalance 查询钱包账户实时余额（GetWalletAccountSingle 别名，兼容旧调用）。
func (g *Gzy) GetBalance(req GetBalanceRequest) (*BalanceResponse, error) {
	return g.GetWalletAccountSingle(req)
}

func (g *Gzy) GetBalanceHistory(req GetBalanceHistoryRequest) (*BalanceHistoryResponse, error) {
	baseURL := g.BaseURL + pathBalanceHistory

	params := url.Values{}
	if req.FromCreated != 0 {
		params.Add("from_created", fmt.Sprintf("%d", req.FromCreated))
	}
	if req.ToCreated != 0 {
		params.Add("to_created", fmt.Sprintf("%d", req.ToCreated))
	}
	if req.Currency != "" {
		params.Add("currency", req.Currency)
	}
	if req.AccountType > 0 {
		params.Add("account_type", fmt.Sprintf("%d", req.AccountType))
	}
	if req.PageSize > 0 {
		params.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		params.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	reqURL := baseURL + "?" + params.Encode()

	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out BalanceHistoryResponse
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (g *Gzy) ApplyCardHolder(req CardHolderApplyRequest) (*CardHolderApplyResponse, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gzy addCardholder: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathAddCardholder
	hreq, err := g.newRequest("POST", reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy addCardholder: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env addCardholderV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy addCardholder: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy addCardholder: 响应缺少 data")
	}
	var out CardHolderApplyResponse
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, fmt.Errorf("gzy addCardholder: decode data: %w", err)
	}
	return &out, nil
}

func (g *Gzy) EditCardHolder(req CardHolderEditRequest) (*CardHolderEditResponse, error) {
	if strings.TrimSpace(req.CardholderID) == "" {
		return nil, fmt.Errorf("gzy editCardholder: cardholderId 不能为空")
	}
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gzy editCardholder: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathEditCardholder
	hreq, err := g.newRequest("POST", reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy editCardholder: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env editCardholderV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy editCardholder: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy editCardholder: 响应缺少 data")
	}
	var out CardHolderEditResponse
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, fmt.Errorf("gzy editCardholder: decode data: %w", err)
	}
	return &out, nil
}

// CreateMatrixAccount POST /matrix/openApi/v4/createMatrixAccount 单一会员下创建 matrix 账户。
func (g *Gzy) CreateMatrixAccount(req CreateMatrixAccountRequest) (*CreateMatrixAccountResponse, error) {
	name := strings.TrimSpace(req.MatrixAccountName)
	if name == "" {
		return nil, fmt.Errorf("gzy createMatrixAccount: matrixAccountName 不能为空")
	}
	bodyStruct := CreateMatrixAccountRequest{MatrixAccountName: name}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("gzy createMatrixAccount: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathCreateMatrixAccount
	hreq, err := g.newRequest(http.MethodPost, reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy createMatrixAccount: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env createMatrixAccountV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy createMatrixAccount: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy createMatrixAccount: 响应缺少 data")
	}
	var out CreateMatrixAccountResponse
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, fmt.Errorf("gzy createMatrixAccount: decode data: %w", err)
	}
	if strings.TrimSpace(out.MatrixAccount) == "" {
		return nil, fmt.Errorf("gzy createMatrixAccount: 响应缺少 data.matrixAccount")
	}
	return &out, nil
}

// addCardholderV4Envelope POST /vcc/openApi/v4/addCardholder 应答外层。
type addCardholderV4Envelope struct {
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
	Method string          `json:"method,omitempty"`
	Path   string          `json:"path,omitempty"`
}

// editCardholderV4Envelope POST /vcc/openApi/v4/editCardholder 应答外层。
type editCardholderV4Envelope struct {
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
	Method string          `json:"method,omitempty"`
	Path   string          `json:"path,omitempty"`
}

// createMatrixAccountV4Envelope POST /matrix/openApi/v4/createMatrixAccount 应答外层。
type createMatrixAccountV4Envelope struct {
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
	Method string          `json:"method,omitempty"`
	Path   string          `json:"path,omitempty"`
}

func (g *Gzy) GetCardHolderDetail(req GetCardHolderDetailRequest) (*CardholderDetail, error) {
	baseURL := g.BaseURL + pathCardholderDetail
	params := url.Values{}
	params.Add("partner_holder_id", req.PartnerHolderID)
	params.Add("card_holder_id", req.CardHolderID)

	reqURL := baseURL + "?" + params.Encode()

	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out CardholderDetail
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// PagingVccCard GET /vcc/openApi/v4/pagingVccCard 查询卡列表。
func (g *Gzy) PagingVccCard(req PagingVccCardRequest) (*CardsPageResponse, error) {
	q := buildPagingVccCardQuery(req)
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathPagingVccCard
	if enc := q.Encode(); enc != "" {
		reqURL += "?" + enc
	}
	hreq, err := g.newRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy pagingVccCard: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env pagingVccCardV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy pagingVccCard: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	var rows []CardPageItem
	if len(env.Data) > 0 && string(bytes.TrimSpace(env.Data)) != "null" {
		if err := json.Unmarshal(env.Data, &rows); err != nil {
			return nil, fmt.Errorf("gzy pagingVccCard: decode data: %w", err)
		}
	}
	out := &CardsPageResponse{
		Numbers:   env.Numbers,
		PageIndex: env.PageIndex,
		PageSize:  env.PageSize,
		Total:     env.Total,
		List:      rows,
	}
	if out.Numbers == 0 && len(rows) > 0 {
		out.Numbers = int32(len(rows))
	}
	if out.PageSize > 0 && out.Total > 0 {
		out.Pages = int((out.Total + out.PageSize - 1) / out.PageSize)
	}
	return out, nil
}

func (g *Gzy) GetCardHoldersPage(req GetCardHoldersPageRequest) (*CardholdersPageResponse, error) {
	q := buildPagingVccCardholderQuery(req)
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathPagingVccCardholder
	if enc := q.Encode(); enc != "" {
		reqURL += "?" + enc
	}
	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy pagingVccCardholder: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env pagingVccCardholderV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy pagingVccCardholder: decode: %w", err)
	}
	if strings.TrimSpace(env.Code) != "0000" {
		return nil, gzyAPIFailure(strings.TrimSpace(env.Code), strings.TrimSpace(env.Msg))
	}
	var rows []CardholderPageItem
	if len(env.Data) > 0 && string(bytes.TrimSpace(env.Data)) != "null" {
		if err := json.Unmarshal(env.Data, &rows); err != nil {
			return nil, fmt.Errorf("gzy pagingVccCardholder: decode data: %w", err)
		}
	}
	out := &CardholdersPageResponse{
		Numbers:   env.Numbers,
		PageIndex: env.PageIndex,
		PageSize:  env.PageSize,
		Total:     env.Total,
		List:      rows,
	}
	if out.Numbers == 0 && len(rows) > 0 {
		out.Numbers = int32(len(rows))
	}
	if out.PageSize > 0 && out.Total > 0 {
		out.Pages = int((out.Total + out.PageSize - 1) / out.PageSize)
	}
	return out, nil
}

// pagingVccCardV4Envelope GET /vcc/openApi/v4/pagingVccCard 应答外层。
type pagingVccCardV4Envelope struct {
	Code      string          `json:"code"`
	Msg       string          `json:"msg"`
	Data      json.RawMessage `json:"data"`
	Numbers   int32           `json:"numbers"`
	PageIndex int64           `json:"pageIndex"`
	PageSize  int64           `json:"pageSize"`
	Total     int64           `json:"total"`
	Method    string          `json:"method,omitempty"`
	Path      string          `json:"path,omitempty"`
}

// pagingVccCardholderV4Envelope GET /vcc/openApi/v4/pagingVccCardholder 应答外层。
type pagingVccCardholderV4Envelope struct {
	Code      string          `json:"code"`
	Msg       string          `json:"msg"`
	Data      json.RawMessage `json:"data"`
	Numbers   int32           `json:"numbers"`
	PageIndex int64           `json:"pageIndex"`
	PageSize  int64           `json:"pageSize"`
	Total     int64           `json:"total"`
	Method    string          `json:"method,omitempty"`
	Path      string          `json:"path,omitempty"`
}

func (g *Gzy) CreateCard(req CreateCardRequest) (*CreateCardResponse, error) {
	if err := validateCreateCardOpenCardV4(&req); err != nil {
		return nil, err
	}
	wire, err := createCardRequestToOpenCardV4(&req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(wire)
	if err != nil {
		return nil, fmt.Errorf("gzy CreateCard: marshal body: %w", err)
	}

	reqURL := strings.TrimRight(g.BaseURL, "/") + pathOpenCard

	hreq, err := g.newRequest("POST", reqURL, jsonBytes)
	if err != nil {
		return nil, err
	}
	logGzyOpenCardRequest(hreq, jsonBytes)

	resp, err := g.client.Do(hreq)
	if err != nil {
		global.GVA_LOG.Error("gzy openCard HTTP transport failed", zap.Error(err))
		return nil, fmt.Errorf("gzy CreateCard: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := readBody(resp)
	logGzyOpenCardResponse(resp.StatusCode, respBody, err)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	var parsed openCardV4Envelope
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("gzy CreateCard: decode response: %w", err)
	}
	out := &CreateCardResponse{
		PartnerOrderID: req.PartnerOrderID,
	}
	if parsed.Data != nil {
		out.RequestID = strings.TrimSpace(parsed.Data.RequestID)
		out.Status = strings.TrimSpace(parsed.Data.Status)
		out.CardDetail = parsed.Data.CardDetail
		if parsed.Data.CardDetail != nil {
			out.CardID = strings.TrimSpace(parsed.Data.CardDetail.CardID)
		}
	}
	logGzyOpenCardParsed(req.PartnerOrderID, &parsed, out)
	if parsed.Code != "0000" {
		return nil, gzyAPIFailure(parsed.Code, parsed.Msg)
	}
	if strings.EqualFold(out.Status, "failed") {
		msg := strings.TrimSpace(parsed.Msg)
		if msg == "" {
			msg = "openCard data.status=failed"
		}
		code := strings.TrimSpace(parsed.Code)
		if code != "" && !strings.EqualFold(code, "0000") {
			return nil, gzyAPIFailure(code, msg)
		}
		return nil, fmt.Errorf("gzy API: message=%s", msg)
	}
	if out.CardID == "" {
		return nil, fmt.Errorf("gzy CreateCard: 响应缺少 data.cardDetail.cardId")
	}
	return out, nil
}

func (g *Gzy) CardFrozen(req CardFrozenRequest) (*string, error) {
	return g.freezeCardV4(req.CardID, req.PartnerOrderID, "freeze")
}

func (g *Gzy) CardUnFrozen(req CardUnFrozenRequest) (*string, error) {
	return g.freezeCardV4(req.CardID, req.PartnerOrderID, "unfreeze")
}

func (g *Gzy) freezeCardV4(cardID, requestID, status string) (*string, error) {
	if strings.TrimSpace(cardID) == "" {
		return nil, fmt.Errorf("gzy freezeCard: cardId 必填")
	}
	if strings.TrimSpace(requestID) == "" {
		return nil, fmt.Errorf("gzy freezeCard: requestId 必填")
	}
	if status != "freeze" && status != "unfreeze" {
		return nil, fmt.Errorf("gzy freezeCard: status 须为 freeze 或 unfreeze")
	}
	bodyStruct := struct {
		CardID    string `json:"cardId"`
		RequestID string `json:"requestId"`
		Status    string `json:"status"`
	}{
		CardID:    strings.TrimSpace(cardID),
		RequestID: strings.TrimSpace(requestID),
		Status:    status,
	}
	jsonBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("gzy freezeCard: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathFreezeCard
	hreq, err := g.newRequest("POST", reqURL, jsonBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy freezeCard: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, respBody); err != nil {
		return nil, err
	}
	var env freezeCardV4Envelope
	if err := json.Unmarshal(respBody, &env); err != nil {
		return nil, fmt.Errorf("gzy freezeCard: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	ok := ""
	return &ok, nil
}

// getCardDetailV4Envelope GET /vcc/openApi/v4/getCardDetail 应答外层。
type getCardDetailV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func photonCardTypeToCardModel(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "share":
		return "SHARE"
	case "recharge":
		return "CARD"
	default:
		return strings.TrimSpace(t)
	}
}

func vccCardInfoToQueryCardDetailResponse(req QueryCardDetailRequest, info *GetCardDetailV4CardInfo) *QueryCardDetailResponse {
	if info == nil {
		return &QueryCardDetailResponse{PartnerOrderID: req.PartnerOrderID, CardID: strings.TrimSpace(req.CardID)}
	}
	cardNo := strings.TrimSpace(info.CardNo)
	if cardNo == "" {
		cardNo = strings.TrimSpace(info.MaskCardNo)
	}
	bal := info.CardBalance
	if bal.IsZero() && !info.AvailableTransactionLimit.IsZero() {
		bal = info.AvailableTransactionLimit
	}
	return &QueryCardDetailResponse{
		PartnerOrderID:   strings.TrimSpace(req.PartnerOrderID),
		CardID:           strings.TrimSpace(info.CardID),
		CardNumber:       cardNo,
		CVV:              strings.TrimSpace(info.CVV),
		Expiry:           strings.TrimSpace(info.ExpirationDate),
		Currency:         strings.TrimSpace(info.CardCurrency),
		ActiveDate:       strings.TrimSpace(info.CreatedAt),
		InactiveDate:     "",
		CardBrand:        NormalizeCardScheme(info.CardScheme),
		CardModel:        photonCardTypeToCardModel(info.CardType),
		CardLevel:        "",
		CardStatus:       PhotonCardStatusToSystem(info.CardStatus),
		AvailableBalance: bal,
		TotalAuthLimit:   info.TotalTransactionLimit,
		UsedAuthLimit:    decimal.Zero,
		PrimaryCardID:    "",
	}
}

func (g *Gzy) QueryCardDetail(req QueryCardDetailRequest) (*QueryCardDetailResponse, error) {
	cardID := strings.TrimSpace(req.CardID)
	if cardID == "" {
		return nil, fmt.Errorf("gzy QueryCardDetail: cardId 必填")
	}
	q := url.Values{}
	q.Set("cardId", cardID)
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathGetCardDetail + "?" + q.Encode()
	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gzy QueryCardDetail: create request: %w", err)
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy QueryCardDetail: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env getCardDetailV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy QueryCardDetail: decode envelope: %w", err)
	}
	code := strings.TrimSpace(env.Code)
	if code != "" && code != "0000" {
		return nil, gzyAPIFailure(code, strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy QueryCardDetail: 响应缺少 data")
	}
	var info GetCardDetailV4CardInfo
	if err := json.Unmarshal(env.Data, &info); err != nil {
		return nil, fmt.Errorf("gzy QueryCardDetail: decode data: %w", err)
	}
	return vccCardInfoToQueryCardDetailResponse(req, &info), nil
}

// getCvvV4Envelope GET /vcc/openApi/v4/getCvv 应答外层。
type getCvvV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// GetCvv GET /vcc/openApi/v4/getCvv?cardId= 查询卡 CVV 与有效期。
func (g *Gzy) GetCvv(req GetCvvRequest) (*VccCvvInfo, error) {
	cardID := strings.TrimSpace(req.CardID)
	if cardID == "" {
		return nil, fmt.Errorf("gzy GetCvv: cardId 必填")
	}
	q := url.Values{}
	q.Set("cardId", cardID)
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathGetCvv + "?" + q.Encode()
	hreq, err := g.newRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gzy GetCvv: create request: %w", err)
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy GetCvv: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env getCvvV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy GetCvv: decode envelope: %w", err)
	}
	code := strings.TrimSpace(env.Code)
	if code != "" && code != "0000" {
		return nil, gzyAPIFailure(code, strings.TrimSpace(env.Msg))
	}
	if len(bytes.TrimSpace(env.Data)) == 0 || string(bytes.TrimSpace(env.Data)) == "null" {
		return nil, fmt.Errorf("gzy GetCvv: 响应缺少 data")
	}
	var out VccCvvInfo
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return nil, fmt.Errorf("gzy GetCvv: decode data: %w", err)
	}
	return &out, nil
}

// getCardBinV4Envelope GET /vcc/openApi/v4/getCardBin 应答外层。
type getCardBinV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// ListCardBin GET /vcc/openApi/v4/getCardBin 查询可用卡 BIN 列表。
func (g *Gzy) ListCardBin() ([]CardBinItem, error) {
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathGetCardBin
	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gzy ListCardBin: create request: %w", err)
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy ListCardBin: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env getCardBinV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy ListCardBin: decode envelope: %w", err)
	}
	code := strings.TrimSpace(env.Code)
	if code != "" && code != "0000" {
		return nil, gzyAPIFailure(code, strings.TrimSpace(env.Msg))
	}
	var rows []CardBinItem
	if len(env.Data) > 0 && string(bytes.TrimSpace(env.Data)) != "null" {
		if err := json.Unmarshal(env.Data, &rows); err != nil {
			return nil, fmt.Errorf("gzy ListCardBin: decode data: %w", err)
		}
	}
	return rows, nil
}

func (g *Gzy) ChangeSubAuthLimit(req ChangeSubAuthLimitRequest) (*string, error) {
	baseURL := g.BaseURL + pathCardOpChangeLimit

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := g.newRequest("POST", baseURL, jsonBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out string
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// --- GET /vcc/openApi/v4/preRecharge；POST /vcc/openApi/v4/recharge（两步转入）---

type preRechargeV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type preRechargeV4Data struct {
	AccountID              string          `json:"accountId"`
	ArrivalAmount          decimal.Decimal `json:"arrivalAmount"`
	ArrivalAmountCurrency  string          `json:"arrivalAmountCurrency"`
	EffectiveQuotationTime int64           `json:"effectiveQuotationTime"`
	ExchangeRate           decimal.Decimal `json:"exchangeRate"`
	QuotedAt               string          `json:"quotedAt"`
	QuotationTime          string          `json:"quotationTime"`
	RechargeAmount         decimal.Decimal `json:"rechargeAmount"`
	RechargeCurrency       string          `json:"rechargeCurrency"`
	RechargeFee            decimal.Decimal `json:"rechargeFee"`
	RechargeFeeCurrency    string          `json:"rechargeFeeCurrency"`
	RequestID              string          `json:"requestId"`
}

type rechargeV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type rechargeV4Data struct {
	ArrivalAmount         decimal.Decimal `json:"arrivalAmount"`
	ArrivalAmountCurrency string          `json:"arrivalAmountCurrency"`
	CardBalance           decimal.Decimal `json:"cardBalance"`
	CardID                string          `json:"cardId"`
	CreatedAt             string          `json:"createdAt"`
	ExchangeRate          decimal.Decimal `json:"exchangeRate"`
	RechargeAmount        decimal.Decimal `json:"rechargeAmount"`
	RechargeCurrency      string          `json:"rechargeCurrency"`
	RechargeFee           decimal.Decimal `json:"rechargeFee"`
	RechargeFeeCurrency   string          `json:"rechargeFeeCurrency"`
	Status                string          `json:"status"`
	TransactionID         string          `json:"transactionId"`
}

func validatePreRechargeRequest(req *PreRechargeRequest) error {
	if strings.TrimSpace(req.RequestID) == "" {
		return fmt.Errorf("gzy preRecharge: requestId 必填")
	}
	if strings.TrimSpace(req.AccountID) == "" {
		return fmt.Errorf("gzy preRecharge: accountId 必填")
	}
	if strings.TrimSpace(req.CardID) == "" {
		return fmt.Errorf("gzy preRecharge: cardId 必填")
	}
	hasR := req.RechargeAmount != nil && req.RechargeAmount.IsPositive()
	hasA := req.ArrivalAmount != nil && req.ArrivalAmount.IsPositive()
	if hasR == hasA {
		return fmt.Errorf("gzy preRecharge: rechargeAmount 与 arrivalAmount 须且只能填其一（正数）")
	}
	return nil
}

func (g *Gzy) PreRecharge(req PreRechargeRequest) (*PreRechargeResponse, error) {
	if err := validatePreRechargeRequest(&req); err != nil {
		return nil, err
	}
	params := url.Values{}
	if mid := strings.TrimSpace(req.MemberID); mid != "" {
		params.Set("memberId", mid)
	}
	params.Set("requestId", strings.TrimSpace(req.RequestID))
	params.Set("accountId", strings.TrimSpace(req.AccountID))
	params.Set("cardId", strings.TrimSpace(req.CardID))
	if req.RechargeAmount != nil && req.RechargeAmount.IsPositive() {
		params.Set("rechargeAmount", req.RechargeAmount.String())
	}
	if req.ArrivalAmount != nil && req.ArrivalAmount.IsPositive() {
		params.Set("arrivalAmount", req.ArrivalAmount.String())
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathPreRecharge + "?" + params.Encode()
	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy preRecharge: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env preRechargeV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy preRecharge: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	var d preRechargeV4Data
	if len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, &d); err != nil {
			return nil, fmt.Errorf("gzy preRecharge: decode data: %w", err)
		}
	}
	quoted := strings.TrimSpace(d.QuotedAt)
	if quoted == "" {
		quoted = strings.TrimSpace(d.QuotationTime)
	}
	qid := strings.TrimSpace(d.RequestID)
	if qid == "" {
		return nil, fmt.Errorf("gzy preRecharge: 响应缺少 data.requestId，无法转入下单")
	}
	return &PreRechargeResponse{
		AccountID:              strings.TrimSpace(d.AccountID),
		ArrivalAmount:          d.ArrivalAmount,
		ArrivalAmountCurrency:  strings.TrimSpace(d.ArrivalAmountCurrency),
		EffectiveQuotationTime: d.EffectiveQuotationTime,
		ExchangeRate:           d.ExchangeRate,
		QuotedAt:               quoted,
		RechargeAmount:         d.RechargeAmount,
		RechargeCurrency:       strings.TrimSpace(d.RechargeCurrency),
		RechargeFee:            d.RechargeFee,
		RechargeFeeCurrency:    strings.TrimSpace(d.RechargeFeeCurrency),
		QuotationRequestID:     qid,
	}, nil
}

func (g *Gzy) RechargeCard(req RechargeCommitRequest) (*RechargeResponse, error) {
	rid := strings.TrimSpace(req.RequestID)
	if rid == "" {
		return nil, fmt.Errorf("gzy recharge: requestId 必填（换汇询价返回的 data.requestId）")
	}
	bodyStruct := struct {
		MemberID  string `json:"memberId,omitempty"`
		RequestID string `json:"requestId"`
	}{
		RequestID: rid,
	}
	if mid := strings.TrimSpace(req.MemberID); mid != "" {
		bodyStruct.MemberID = mid
	}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("gzy recharge: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathRecharge
	hreq, err := g.newRequest("POST", reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy recharge: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env rechargeV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy recharge: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	var d rechargeV4Data
	if len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, &d); err != nil {
			return nil, fmt.Errorf("gzy recharge: decode data: %w", err)
		}
	}
	out := &RechargeResponse{
		PartnerOrderID:        rid,
		Status:                strings.TrimSpace(d.Status),
		CardID:                strings.TrimSpace(d.CardID),
		TransactionID:         strings.TrimSpace(d.TransactionID),
		ArrivalAmount:         d.ArrivalAmount,
		ArrivalAmountCurrency: strings.TrimSpace(d.ArrivalAmountCurrency),
		CardBalance:           d.CardBalance,
		CreatedAt:             strings.TrimSpace(d.CreatedAt),
		ExchangeRate:          d.ExchangeRate,
		RechargeAmount:        d.RechargeAmount,
		RechargeCurrency:      strings.TrimSpace(d.RechargeCurrency),
		RechargeFee:           d.RechargeFee,
		RechargeFeeCurrency:   strings.TrimSpace(d.RechargeFeeCurrency),
	}
	if strings.EqualFold(out.Status, "failed") {
		msg := strings.TrimSpace(env.Msg)
		if msg == "" {
			msg = "recharge status=failed"
		}
		return out, gzyAPIFailure(env.Code, msg)
	}
	return out, nil
}

func (g *Gzy) WithdrawFromCard(req WithdrawRequest) (*WithdrawResponse, error) {
	cardID := strings.TrimSpace(req.CardID)
	requestID := strings.TrimSpace(req.PartnerOrderID)
	if cardID == "" {
		return nil, fmt.Errorf("gzy rechargeReturn: cardId 必填")
	}
	if requestID == "" {
		return nil, fmt.Errorf("gzy rechargeReturn: requestId 必填")
	}
	if !req.Amount.IsPositive() {
		return nil, fmt.Errorf("gzy rechargeReturn: returnAmount 须为正数")
	}
	bodyStruct := struct {
		CardID       string          `json:"cardId"`
		RequestID    string          `json:"requestId"`
		ReturnAmount decimal.Decimal `json:"returnAmount"`
	}{
		CardID:       cardID,
		RequestID:    requestID,
		ReturnAmount: req.Amount,
	}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("gzy rechargeReturn: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathRechargeReturn
	hreq, err := g.newRequest("POST", reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy rechargeReturn: %w", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env rechargeReturnV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy rechargeReturn: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	var d rechargeReturnV4Data
	if len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, &d); err != nil {
			return nil, fmt.Errorf("gzy rechargeReturn: decode data: %w", err)
		}
	}
	outCardID := strings.TrimSpace(d.CardID)
	if outCardID == "" {
		outCardID = cardID
	}
	out := &WithdrawResponse{
		PartnerOrderID:  requestID,
		CardID:          outCardID,
		TransactionID:   strings.TrimSpace(d.TransactionID),
		CreatedAt:       strings.TrimSpace(d.CreatedAt),
		MaskCardNo:      strings.TrimSpace(d.MaskCardNo),
		ArrivalAmount:   d.ArrivalAmount,
		ReturnFeeAmount: d.ReturnFeeAmount,
		CardBalance:     d.CardBalance,
		Status:          strings.TrimSpace(d.Status),
	}
	if strings.EqualFold(out.Status, "failed") {
		msg := strings.TrimSpace(env.Msg)
		if msg == "" {
			msg = "rechargeReturn status=failed"
		}
		return out, gzyAPIFailure(env.Code, msg)
	}
	return out, nil
}

func (g *Gzy) CancelCard(req CancelCardRequest) (*CancelCardResponse, error) {
	cardID := strings.TrimSpace(req.CardID)
	if cardID == "" {
		return nil, fmt.Errorf("gzy cancelCard: cardId 必填")
	}
	bodyStruct := struct {
		CardID string `json:"cardId"`
	}{CardID: cardID}
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("gzy cancelCard: marshal: %w", err)
	}
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathCardOpCancel
	hreq, err := g.newRequest("POST", reqURL, bodyBytes)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy cancelCard: %w", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env cancelCardV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy cancelCard: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	out := &CancelCardResponse{
		PartnerOrderID: strings.TrimSpace(req.PartnerOrderID),
		CardID:         cardID,
	}
	if len(env.Data) > 0 {
		var d cancelCardV4Data
		if err := json.Unmarshal(env.Data, &d); err == nil {
			if tid := strings.TrimSpace(d.TransactionID); tid != "" {
				out.TransactionID = tid
			} else if tid := strings.TrimSpace(d.TxnIDSnake); tid != "" {
				out.TransactionID = tid
			}
			if cid := strings.TrimSpace(d.CardID); cid != "" {
				out.CardID = cid
			} else if cid := strings.TrimSpace(d.CardIDSnake); cid != "" {
				out.CardID = cid
			}
		}
	}
	return out, nil
}

func (g *Gzy) ApplyInbound(req ApplyInboundRequest) (*ApplyInboundResp, error) {
	url := g.BaseURL + pathFundInboundApply

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	hreq, err := g.newRequest("POST", url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out ApplyInboundResp
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (g *Gzy) GetInboundDetail(req GetInboundRequest) (*InboundDetail, error) {
	baseURL := g.BaseURL + pathQueryCardTransactions
	queryParams := make(url.Values)
	queryParams.Add("partner_order_id", req.PartnerOrdeID)
	queryParams.Add("order_id", req.OrderId)

	url := baseURL + "?" + queryParams.Encode()

	hreq, err := g.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out InboundDetail
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (g *Gzy) QueryCardTransactions(req QueryCardTransactionsRequest) (*QueryCardTransactionsResponse, error) {
	q := buildQueryCardTransactionsV4Query(req)
	reqURL := strings.TrimRight(g.BaseURL, "/") + pathPagingVccTradeOrder
	if enc := q.Encode(); enc != "" {
		reqURL += "?" + enc
	}
	hreq, err := g.newRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("gzy queryCardTransactions: %w", err)
	}
	defer resp.Body.Close()
	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var env queryCardTransactionsV4Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("gzy queryCardTransactions: decode: %w", err)
	}
	if env.Code != "0000" {
		return nil, gzyAPIFailure(env.Code, env.Msg)
	}
	var rows []VccTradeOrderResp
	if len(env.Data) > 0 && string(bytes.TrimSpace(env.Data)) != "null" {
		if err := json.Unmarshal(env.Data, &rows); err != nil {
			return nil, fmt.Errorf("gzy queryCardTransactions: decode data: %w", err)
		}
	}
	out := &QueryCardTransactionsResponse{
		Numbers:   env.Numbers,
		PageIndex: env.PageIndex,
		PageSize:  env.PageSize,
		Total:     env.Total,
		List:      rows,
	}
	if out.Numbers == 0 && len(rows) > 0 {
		out.Numbers = int32(len(rows))
	}
	if out.PageSize > 0 && out.Total > 0 {
		out.Pages = int((out.Total + out.PageSize - 1) / out.PageSize)
	}
	return out, nil
}

func (g *Gzy) ListInboundDetails(req InboundQueryRequest) (*InboundListResponse, error) {
	baseURL := g.BaseURL + pathFundInboundList
	queryParams := make(url.Values)
	if req.State != "" {
		queryParams.Add("state", req.State)
	}
	if req.ChainName != "" {
		queryParams.Add("chain_name", req.ChainName)
	}
	if req.OrderType != "" {
		queryParams.Add("order_type", req.OrderType)
	}
	if req.BeginCreateTime != "" {
		queryParams.Add("begin_create_time", req.BeginCreateTime)
	}
	if req.EndCreateTime != "" {
		queryParams.Add("end_create_time", req.EndCreateTime)
	}
	if req.BeginFinishTime != "" {
		queryParams.Add("begin_finish_time", req.BeginFinishTime)
	}
	if req.EndFinishTime != "" {
		queryParams.Add("end_finish_time", req.EndFinishTime)
	}
	if req.PageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", req.PageSize))
	}
	if req.PageNo > 0 {
		queryParams.Add("page_no", fmt.Sprintf("%d", req.PageNo))
	}

	url := baseURL + "?" + queryParams.Encode()

	hreq, err := g.newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := g.client
	resp, err := client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}
	if err := httpPhotonOrBodyError(resp.StatusCode, body); err != nil {
		return nil, err
	}
	var out InboundListResponse
	if err := decodePhotonJSON(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// --- GET /vcc/openApi/v4/pagingVccTradeOrder（QueryCardTransactions） ---

type queryCardTransactionsV4Envelope struct {
	Code      string          `json:"code"`
	Msg       string          `json:"msg"`
	Data      json.RawMessage `json:"data"`
	Numbers   int32           `json:"numbers"`
	PageIndex int64           `json:"pageIndex"`
	PageSize  int64           `json:"pageSize"`
	Total     int64           `json:"total"`
	Method    string          `json:"method,omitempty"`
	Path      string          `json:"path,omitempty"`
}

// --- POST /vcc/openApi/v4/rechargeReturn（WithdrawFromCard） ---

type rechargeReturnV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type rechargeReturnV4Data struct {
	TransactionID   string          `json:"transactionId"`
	CardID          string          `json:"cardId"`
	CreatedAt       string          `json:"createdAt"`
	MaskCardNo      string          `json:"maskCardNo"`
	ArrivalAmount   decimal.Decimal `json:"arrivalAmount"`
	ReturnFeeAmount decimal.Decimal `json:"returnFeeAmount"`
	CardBalance     decimal.Decimal `json:"cardBalance"`
	Status          string          `json:"status"`
}

// --- POST /vcc/openApi/v4/cancelCard ---

type cancelCardV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type cancelCardV4Data struct {
	CardID        string `json:"cardId"`
	TransactionID string `json:"transactionId"`
	CardIDSnake   string `json:"card_id"`
	TxnIDSnake    string `json:"transaction_id"`
}

// --- POST /vcc/openApi/v4/freezeCard（CardFrozen / CardUnFrozen） ---

type freezeCardV4Envelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// --- POST /vcc/openApi/v4/openCard（由 CreateCard 调用） ---

type openCardV4Wire struct {
	RequestID            string           `json:"requestId"`
	AccountID            string           `json:"accountId"` // 转入所用币种的光子易账户 ID（wallet/account/single 的 accountNo）
	CardBin              string           `json:"cardBin"`
	CardCurrency         string           `json:"cardCurrency"`
	CardType             string           `json:"cardType"`
	CardholderID         string           `json:"cardholderId,omitempty"`
	ArrivalAmount        *decimal.Decimal `json:"arrivalAmount,omitempty"` // 希望到账金额（与 rechargeAmount 择一；openCard 使用到账金额）
	TransactionLimit     *decimal.Decimal `json:"transactionLimit,omitempty"`
	TransactionLimitType string           `json:"transactionLimitType,omitempty"`
}

type openCardV4Envelope struct {
	Code string                  `json:"code"`
	Msg  string                  `json:"msg"`
	Data *OpenCardV4ResponseData `json:"data"`
}

func validateCreateCardOpenCardV4(req *CreateCardRequest) error {
	if strings.TrimSpace(req.PartnerOrderID) == "" {
		return fmt.Errorf("gzy CreateCard: partner_order_id 必填，将作为 requestId")
	}
	if strings.TrimSpace(req.CardBin) == "" {
		return fmt.Errorf("gzy CreateCard: card_bin 必填（Photon cardBin，非 card_bin_id）")
	}
	if strings.TrimSpace(req.AccountCurrency) == "" {
		return fmt.Errorf("gzy CreateCard: account_currency 必填，将作为 cardCurrency")
	}
	if strings.TrimSpace(req.CardModel) == "" {
		return fmt.Errorf("gzy CreateCard: card_model 必填（SHARE|CARD 对应 share|recharge）")
	}
	if resolveOpenCardAccountID(req) == "" {
		return fmt.Errorf("gzy CreateCard: accountId 必填")
	}
	return nil
}

// ResolveAccountID 解析光子易账户 ID：请求 account_id > 配置 gzy.account-id > DefaultGzyAccountID。
func ResolveAccountID(override string) string {
	if s := strings.TrimSpace(override); s != "" {
		return s
	}
	if s := strings.TrimSpace(global.GVA_CONFIG.Gzy.AccountID); s != "" {
		return s
	}
	return DefaultGzyAccountID
}

func resolveOpenCardAccountID(req *CreateCardRequest) string {
	if req != nil {
		return ResolveAccountID(req.AccountID)
	}
	return ResolveAccountID("")
}

func createCardRequestToOpenCardV4(req *CreateCardRequest) (*openCardV4Wire, error) {
	w := &openCardV4Wire{
		RequestID:    strings.TrimSpace(req.PartnerOrderID),
		AccountID:    resolveOpenCardAccountID(req),
		CardBin:      strings.TrimSpace(req.CardBin),
		CardCurrency: strings.TrimSpace(req.AccountCurrency),
		CardType:     cardModelToPhotonV4CardType(req.CardModel),
		CardholderID: strings.TrimSpace(req.CardHolderID),
	}
	if a := strings.TrimSpace(req.Amount); a != "" && a != "0" {
		d, err := decimal.NewFromString(a)
		if err != nil {
			return nil, fmt.Errorf("gzy CreateCard: amount 非法: %w", err)
		}
		if !d.IsZero() {
			w.ArrivalAmount = &d
		}
	}
	if strings.TrimSpace(req.PrimaryCardID) != "" {
		if strings.EqualFold(strings.TrimSpace(req.AuthLimitFlag), "Y") && strings.TrimSpace(req.TotalAuthLimit) != "" {
			lim, err := decimal.NewFromString(strings.TrimSpace(req.TotalAuthLimit))
			if err != nil {
				return nil, fmt.Errorf("gzy CreateCard: total_auth_limit 非法: %w", err)
			}
			w.TransactionLimitType = "limited"
			w.TransactionLimit = &lim
		}
	}
	return w, nil
}

func cardModelToPhotonV4CardType(model string) string {
	switch strings.ToUpper(strings.TrimSpace(model)) {
	case "SHARE":
		return "share"
	case "CARD":
		return "recharge"
	default:
		return strings.ToLower(strings.TrimSpace(model))
	}
}
