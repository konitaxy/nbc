package captcha

import (
	"log"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func GenerateTOTP(email string) (string, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Newbeecard.com",
		AccountName: email,
	})
	if err != nil {
		log.Fatal(err)
	}

	return key.Secret(), key.URL()
}

func VerifyTOTP(secret string, userCode string) bool {
	valid := totp.Validate(userCode, secret)
	return valid
}
func GenerateQRCode(url string) (img []byte, err error) {
	return qrcode.Encode(url, qrcode.Medium, 256)
}
