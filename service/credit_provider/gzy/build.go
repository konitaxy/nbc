package gzy

// 本文件集中存放请求构造类辅助函数（query 拼接、X-PD-SIGN 等）。

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitlab.com/ucard/global"
)

const defaultGzySignPrivateKeyPath = "resource/keys/rsa_private.pem"

func resolvedGzySignPrivateKeyPath() string {
	if p := strings.TrimSpace(global.GVA_CONFIG.Gzy.SignPrivateKeyPath); p != "" {
		return p
	}
	return defaultGzySignPrivateKeyPath
}

// hasGzyRSASigningMaterial 是否具备 RSA 加签材料（内联 PEM 或私钥文件存在）。
func hasGzyRSASigningMaterial() bool {
	if strings.TrimSpace(global.GVA_CONFIG.Gzy.SignPrivateKey) != "" {
		return true
	}
	p := resolvedGzySignPrivateKeyPath()
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

// loadGzySignPrivateKeyPEM 优先使用配置内联 PEM，否则读取 sign-private-key-path（空则 resource/keys/rsa_private.pem）。
func loadGzySignPrivateKeyPEM() (string, error) {
	if s := strings.TrimSpace(global.GVA_CONFIG.Gzy.SignPrivateKey); s != "" {
		return s, nil
	}
	path := resolvedGzySignPrivateKeyPath()
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取 gzy 私钥文件 %q: %w", path, err)
	}
	return string(b), nil
}

// buildWalletAccountSingleQuery 构造 GET /wallet/openApi/v4/account/single 的 query。
// 传 accountNo 时仅带 accountNo（及可选 memberId）；否则 currency 必填，accountType 默认 FT10001。
func buildWalletAccountSingleQuery(req GetBalanceRequest) (url.Values, error) {
	v := url.Values{}
	accNo := strings.TrimSpace(req.AccountNo)
	if accNo != "" {
		v.Set("accountNo", accNo)
		if mid := strings.TrimSpace(req.MemberID); mid != "" {
			v.Set("memberId", mid)
		}
		return v, nil
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, fmt.Errorf("currency 必填（ISO4217）；或传入 accountNo 精确查询")
	}
	v.Set("currency", strings.TrimSpace(req.Currency))
	at := strings.TrimSpace(req.AccountType)
	if at == "" {
		at = WalletAccountTypeAvailable
	}
	v.Set("accountType", at)
	if mid := strings.TrimSpace(req.MemberID); mid != "" {
		v.Set("memberId", mid)
	}
	if mx := strings.TrimSpace(req.MatrixAccount); mx != "" {
		v.Set("matrixAccount", mx)
	}
	return v, nil
}

// buildPagingVccCardholderQuery 构造 GET /vcc/openApi/v4/pagingVccCardholder 的 query。
func buildPagingVccCardholderQuery(req GetCardHoldersPageRequest) url.Values {
	v := url.Values{}
	if req.PageIndex > 0 {
		v.Set("pageIndex", fmt.Sprintf("%d", req.PageIndex))
	}
	if req.PageSize > 0 {
		v.Set("pageSize", fmt.Sprintf("%d", req.PageSize))
	}
	if s := strings.TrimSpace(req.MemberID); s != "" {
		v.Set("memberId", s)
	}
	if s := strings.TrimSpace(req.MatrixAccount); s != "" {
		v.Set("matrixAccount", s)
	}
	if s := strings.TrimSpace(req.CreatedAtStart); s != "" {
		v.Set("createdAtStart", s)
	}
	if s := strings.TrimSpace(req.CreatedAtEnd); s != "" {
		v.Set("createdAtEnd", s)
	}
	if s := strings.TrimSpace(req.CardholderID); s != "" {
		v.Set("cardholderId", s)
	}
	if s := strings.TrimSpace(req.Status); s != "" {
		v.Set("status", s)
	}
	return v
}

// buildQueryCardTransactionsV4Query 构造 GET /vcc/openApi/v4/pagingVccTradeOrder 的 query。
func buildQueryCardTransactionsV4Query(req QueryCardTransactionsRequest) url.Values {
	v := url.Values{}
	if req.PageIndex > 0 {
		v.Set("pageIndex", fmt.Sprintf("%d", req.PageIndex))
	}
	if req.PageSize > 0 {
		v.Set("pageSize", fmt.Sprintf("%d", req.PageSize))
	}
	if s := strings.TrimSpace(req.MemberID); s != "" {
		v.Set("memberId", s)
	}
	if s := strings.TrimSpace(req.MatrixAccount); s != "" {
		v.Set("matrixAccount", s)
	}
	if s := strings.TrimSpace(req.CreatedAtStart); s != "" {
		v.Set("createdAtStart", s)
	}
	if s := strings.TrimSpace(req.CreatedAtEnd); s != "" {
		v.Set("createdAtEnd", s)
	}
	if s := strings.TrimSpace(req.CardID); s != "" {
		v.Set("cardId", s)
	}
	if s := strings.TrimSpace(req.CardType); s != "" {
		v.Set("cardType", s)
	}
	if s := strings.TrimSpace(req.CardFormFactor); s != "" {
		v.Set("cardFormFactor", s)
	}
	if s := strings.TrimSpace(req.RequestID); s != "" {
		v.Set("requestId", s)
	}
	if s := strings.TrimSpace(req.TransactionID); s != "" {
		v.Set("transactionId", s)
	}
	if s := strings.TrimSpace(req.TransactionType); s != "" {
		v.Set("transactionType", s)
	}
	if s := strings.TrimSpace(req.Status); s != "" {
		v.Set("status", s)
	}
	if s := strings.TrimSpace(req.Nickname); s != "" {
		v.Set("nickname", s)
	}
	return v
}

// buildXPDSign 生成请求头 X-PD-SIGN。PhotonPay 官方：对 request.body 按 UTF-8 视为原字符串（与 JSON 原文一致），
// MD5 摘要后 RSA PKCS#1 v15 + MD5，再 Base64；body 为空（如 GET 无 body）时无需加签，返回 ("", nil)。
//
// gzy.sign-mode（显式配置时优先生效）:
//   - rsa-md5-pkcs1-base64（别名 md5withrsa）：官方 MD5withRSA + Base64，需 gzy.sign-private-key 内联 PEM 或可读私钥文件（见 sign-private-key-path / 默认 resource/keys/rsa_private.pem）
//   - rsa-sha256-pkcs1-base64：对 body SHA256 摘要后 RSA-PKCS1v15 + Base64（历史联调）
//   - body-hmac-sha256-hex：HMAC-SHA256 + 小写 hex（历史联调）
//
// sign-mode 为空时：若具备 RSA 私钥（内联或文件）则默认 rsa-md5-pkcs1-base64；否则若配置了 sign-secret 则 body-hmac-sha256-hex。
func buildXPDSign(body []byte) (string, error) {
	if len(body) == 0 {
		return "", nil
	}
	cfg := global.GVA_CONFIG.Gzy
	mode := strings.TrimSpace(strings.ToLower(cfg.SignMode))
	if mode == "" {
		if hasGzyRSASigningMaterial() {
			mode = "rsa-md5-pkcs1-base64"
		} else if strings.TrimSpace(cfg.SignSecret) != "" {
			mode = "body-hmac-sha256-hex"
		}
	}

	switch mode {
	case "rsa-md5-pkcs1-base64", "md5withrsa", "rsa-md5":
		pemKey, err := loadGzySignPrivateKeyPEM()
		if err != nil {
			return "", fmt.Errorf("gzy X-PD-SIGN: %w", err)
		}
		s, err := signPDDataMD5RSA(string(body), pemKey)
		if err != nil {
			return "", fmt.Errorf("gzy X-PD-SIGN: %w", err)
		}
		return s, nil

	case "body-hmac-sha256-hex":
		secret := strings.TrimSpace(cfg.SignSecret)
		if secret == "" {
			return "", fmt.Errorf("gzy X-PD-SIGN: body-hmac-sha256-hex 需要配置 gzy.sign-secret")
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		return hex.EncodeToString(mac.Sum(nil)), nil

	case "rsa-sha256-pkcs1-base64":
		pemKey, err := loadGzySignPrivateKeyPEM()
		if err != nil {
			return "", fmt.Errorf("gzy X-PD-SIGN: %w", err)
		}
		priv, err := parseRSAPrivateKeyPEM([]byte(pemKey))
		if err != nil {
			return "", fmt.Errorf("gzy X-PD-SIGN: 解析私钥 PEM: %w", err)
		}
		sum := sha256.Sum256(body)
		sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum[:])
		if err != nil {
			return "", fmt.Errorf("gzy X-PD-SIGN: rsa 签名: %w", err)
		}
		return base64.StdEncoding.EncodeToString(sig), nil

	default:
		return "", fmt.Errorf("gzy X-PD-SIGN: 请在配置中设置 gzy.sign-mode（rsa-md5-pkcs1-base64 | rsa-sha256-pkcs1-base64 | body-hmac-sha256-hex），或与官方文档核对；当前 sign-mode=%q", cfg.SignMode)
	}
}
