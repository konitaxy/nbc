package gzy

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// VerifyXPDSignMD5WithRSA 使用 PhotonPay 公钥对异步通知等报文验签：X-PD-SIGN 为 Base64(MD5withRSA(body))。
// body 为空或未提供签名时返回错误。
func VerifyXPDSignMD5WithRSA(pub *rsa.PublicKey, body []byte, signBase64 string) error {
	if pub == nil {
		return fmt.Errorf("gzy 验签: 公钥为空")
	}
	if len(body) == 0 {
		return fmt.Errorf("gzy 验签: body 为空")
	}
	signBase64 = strings.TrimSpace(signBase64)
	if signBase64 == "" {
		return fmt.Errorf("gzy 验签: X-PD-SIGN 为空")
	}
	sig, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return fmt.Errorf("gzy 验签: Base64 解码 X-PD-SIGN: %w", err)
	}
	h := md5.New()
	_, _ = h.Write(body)
	if err := rsa.VerifyPKCS1v15(pub, crypto.MD5, h.Sum(nil), sig); err != nil {
		return fmt.Errorf("gzy 验签: MD5withRSA 校验失败: %w", err)
	}
	return nil
}

// ParseGzyPubKeyFromConfig 解析 gzy.pub-key：支持 PEM 或 Base64 编码的 PKIX DER（Photon 控制台常见格式）。
func ParseGzyPubKeyFromConfig(raw string) (*rsa.PublicKey, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("gzy.pub-key 未配置")
	}
	if strings.Contains(raw, "BEGIN") {
		return ParseRSAPublicKeyPEM([]byte(raw))
	}
	der, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("gzy.pub-key Base64 解码: %w", err)
	}
	k, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("gzy.pub-key 解析 PKIX: %w", err)
	}
	pk, ok := k.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("gzy.pub-key 非 RSA 公钥")
	}
	return pk, nil
}

// ParseRSAPublicKeyPEM 解析 PKIX 或 PKCS1 格式的 RSA 公钥 PEM（用于 PhotonPay 异步通知验签）。
func ParseRSAPublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("无有效 PEM 块")
	}
	if strings.Contains(block.Type, "PUBLIC KEY") {
		k, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pk, ok := k.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("非 RSA 公钥")
		}
		return pk, nil
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// signPDDataMD5RSA PhotonPay X-PD-SIGN（与官方示例一致）：data 为待签名的原字符串（如 JSON 原文），
// MD5 摘要后做 RSA PKCS#1 v15 + MD5，再 Base64。私钥为 PEM，优先 PKCS#8，兼容 PKCS#1。
// 注：rsa.SignPKCS1v15 在 Go 中需传入 rand.Reader，不可为 nil。
func signPDDataMD5RSA(data string, privateKeyPEM string) (string, error) {
	priv, err := parseRSAPrivateKeyPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("解析私钥 PEM: %w", err)
	}
	h := md5.New()
	if _, err := h.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("MD5 写入数据: %w", err)
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.MD5, h.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("MD5withRSA 签名: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func parseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("无有效 PEM 块")
	}
	if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rk, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("非 RSA 私钥")
		}
		return rk, nil
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("无法解析为 PKCS#8 或 PKCS#1 RSA 私钥")
}
