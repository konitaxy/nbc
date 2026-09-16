package adsvcc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestSignParamsDeterministicOrder(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	params := map[string]string{
		"amount":     "1",
		"product_id": "12",
		"num":        "1",
	}
	sig1, err := SignParams(string(pemBytes), params, "1758081170")
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := SignParams(string(pemBytes), map[string]string{
		"num":        "1",
		"product_id": "12",
		"amount":     "1",
	}, "1758081170")
	if err != nil {
		t.Fatal(err)
	}
	if sig1 != sig2 {
		t.Fatalf("sign not stable under key reorder")
	}
	if sig1 == "" {
		t.Fatal("empty signature")
	}
}
