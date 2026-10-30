package ailink

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestCardBinAmountAcceptsNumber(t *testing.T) {
	raw := []byte(`{"list":[{"binId":12,"bin":"486532","minOpenCardAmount":10,"minRechargeAmount":"1.50"}]}`)
	var out CardBinList
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	if len(out.List) != 1 || out.List[0].MinOpenCardAmount.String() != "10" || out.List[0].MinRechargeAmount.String() != "1.50" {
		t.Fatalf("%+v", out.List)
	}
}

func TestMapBrandMaster(t *testing.T) {
	if got := MapBrand("MASTER"); got != "MasterCard" {
		t.Fatalf("got %q", got)
	}
	if got := MapBrand("MASTERCARD"); got != "MasterCard" {
		t.Fatalf("got %q", got)
	}
}

func TestSignRequestStable(t *testing.T) {
	body := []byte(`{"amount":"1.00"}`)
	a := SignRequest("post", "/openapi/v1/cards/funds/operate", "1", "n", "secret", body)
	b := SignRequest("POST", "/openapi/v1/cards/funds/operate", "1", "n", "secret", body)
	if a == "" || a != b {
		t.Fatalf("sign=%q %q", a, b)
	}
	if SignWebhook("e", "T", "1", "n", "secret", body) == "" {
		t.Fatal("empty webhook sign")
	}
}

func TestDecryptSensitiveRoundTrip(t *testing.T) {
	secret := "whsec"
	plain := "4111111111111111"
	block, err := aes.NewCipher(SensitiveKey(secret))
	if err != nil {
		t.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}
	iv := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		t.Fatal(err)
	}
	sealed := gcm.Seal(nil, iv, []byte(plain), nil)
	raw := append(append([]byte{}, iv...), sealed...)
	enc := base64.RawURLEncoding.EncodeToString(raw)
	got, err := DecryptSensitive(secret, enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != plain {
		t.Fatalf("got %q", got)
	}
}

func TestVerifyWebhook(t *testing.T) {
	body := []byte(`{"eventId":"e1"}`)
	sig := "v1=" + SignWebhook("e1", "WEBHOOK_TEST", "10", "abc", "secret", body)
	if err := VerifyWebhook("secret", body, "e1", "WEBHOOK_TEST", "10", "abc", sig); err != nil {
		t.Fatal(err)
	}
	if err := VerifyWebhook("secret", body, "e1", "WEBHOOK_TEST", "10", "abc", "v1=00"); err == nil {
		t.Fatal("expected mismatch")
	}
}
