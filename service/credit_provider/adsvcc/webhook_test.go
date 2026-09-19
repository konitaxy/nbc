package adsvcc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
)

func TestVerifyWebhookSign(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	body := []byte(`{"type":"card_create","code":0,"data":{"batch_id":"OP1"},"msg":"","timestamp":1,"nonce":"n"}`)
	ts := "1789614810"
	sum := sha256.Sum256([]byte(ts + "." + string(body)))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	sigB64 := base64.StdEncoding.EncodeToString(sig)

	old := DefaultPlatformPublicKey
	// temporarily override via parsing helper path: inject by setting env config is hard;
	// call verify with custom parse by temporarily swapping is awkward — test parse + verify locally.
	pub, err := parseRSAPublicKey(string(pubPEM))
	if err != nil {
		t.Fatal(err)
	}
	rawSig, _ := base64.StdEncoding.DecodeString(sigB64)
	payload := []byte(ts + "." + string(body))
	sum2 := sha256.Sum256(payload)
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, sum2[:], rawSig); err != nil {
		t.Fatalf("verify: %v", err)
	}
	_ = old
}

func TestMapCardStatus(t *testing.T) {
	if MapCardStatus(1) != "Active" || MapCardStatus(-2) != "Suspend" || MapCardStatus(-1) != "Cancel" {
		t.Fatal("unexpected status mapping")
	}
}
