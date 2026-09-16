package adsvcc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

// SignParams 按文档：body 参数 + timestamp（字符串）合并 → key 升序 JSON → RSA-SHA256 → Base64。
func SignParams(privateKeyPEM string, params map[string]string, timestampSec string) (string, error) {
	key, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}
	payload := make(map[string]string, len(params)+1)
	for k, v := range params {
		payload[k] = v
	}
	payload["timestamp"] = timestampSec
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("adsvcc sign marshal: %w", err)
	}
	sum := sha256.Sum256(raw)
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		return "", fmt.Errorf("adsvcc sign: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func parseRSAPrivateKey(pemOrRaw string) (*rsa.PrivateKey, error) {
	s := strings.TrimSpace(pemOrRaw)
	if s == "" {
		return nil, fmt.Errorf("adsvcc: private key empty")
	}
	var der []byte
	if block, _ := pem.Decode([]byte(s)); block != nil {
		der = block.Bytes
	} else {
		b, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(s, "\n", ""))
		if err != nil {
			return nil, fmt.Errorf("adsvcc: private key not PEM/Base64")
		}
		der = b
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("adsvcc: private key is not RSA")
		}
		return rsaKey, nil
	}
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("adsvcc: parse private key failed")
}

func loadPrivateKeyPEM(inline, path string) (string, error) {
	if s := strings.TrimSpace(inline); s != "" {
		return s, nil
	}
	path = strings.TrimSpace(path)
	if path == "" {
		path = "resource/keys/adsvcc_private.pem"
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("adsvcc: read private key: %w", err)
	}
	return string(b), nil
}
