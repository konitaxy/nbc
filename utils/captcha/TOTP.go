package captcha

import (
	"fmt"
	"log"
	"net/url"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

const adminTOTPIssuer = "Newbeecard Admin"

func GenerateAdminTOTP(account string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      adminTOTPIssuer,
		AccountName: account,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func AdminTOTPURL(account, secret string) string {
	label := url.PathEscape(adminTOTPIssuer) + ":" + url.PathEscape(account)
	return fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", label, secret, url.QueryEscape(adminTOTPIssuer))
}

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
