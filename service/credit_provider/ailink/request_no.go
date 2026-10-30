package ailink

import (
	"regexp"
	"strings"
)

var requestNoOK = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{7,63}$`)

// NormalizeRequestNo 满足开卡 requestNo：字母开头，同时含字母和数字，8 到 64 位。
func NormalizeRequestNo(raw string) string {
	var b strings.Builder
	for _, r := range strings.TrimSpace(raw) {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	s := b.String()
	hasLetter := strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasDigit := strings.ContainsAny(s, "0123456789")
	if s == "" || s[0] < 'A' || (s[0] > 'Z' && s[0] < 'a') || s[0] > 'z' || !hasLetter || !hasDigit {
		s = "AL" + s
		if !strings.ContainsAny(s, "0123456789") {
			s += "1"
		}
	}
	if len(s) < 8 {
		s += strings.Repeat("0", 8-len(s))
	}
	if len(s) > 64 {
		s = s[:64]
	}
	if !requestNoOK.MatchString(s) || !strings.ContainsAny(s, "0123456789") {
		return "ALREQ0001"
	}
	return s
}

// MapCardStatus 渠道卡状态 → 业务状态。
func MapCardStatus(st string) string {
	switch strings.ToUpper(strings.TrimSpace(st)) {
	case "ACTIVE":
		return "Active"
	case "FROZEN":
		return "Suspend"
	case "CANCELLED", "CANCELED":
		return "Cancel"
	case "EXPIRED":
		return "Closed"
	case "OPENING", "CANCEL_REQUESTED":
		return "Pending"
	default:
		if strings.TrimSpace(st) == "" {
			return ""
		}
		return "Pending"
	}
}

// MapTxnStatus 交易列表状态。已授权按成功。
func MapTxnStatus(st string) string {
	switch strings.ToUpper(strings.TrimSpace(st)) {
	case "SUCCESS", "AUTHORIZED", "SETTLED", "REFUNDED", "REVERSED",
		"AUTH_SUCCESS", "REFUND_SUCCESS", "REVERSE_SUCCESS", "CORRECTION":
		return "Success"
	case "FAILED", "DECLINED", "DISPUTED", "AUTH_FAILED":
		return "Failure"
	case "PROCESSING":
		return "Pending"
	default:
		return strings.TrimSpace(st)
	}
}

func MapBrand(org string) string {
	switch strings.ToUpper(strings.TrimSpace(org)) {
	case "VISA":
		return "Visa"
	case "MASTER", "MASTERCARD":
		return "MasterCard"
	default:
		return strings.TrimSpace(org)
	}
}
