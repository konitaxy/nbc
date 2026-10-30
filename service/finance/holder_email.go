package finance

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	modelfinance "gitlab.com/ucard/model/finance"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

var randomEmailDomains = []string{"gmail.com", "outlook.com", "yahoo.com", "hotmail.com"}

func lookupUserEmail(iamID, clientID uint, isIAM bool) string {
	if isIAM && iamID > 0 {
		var user client.IAMUser
		if err := global.GVA_DB.Select("email").First(&user, "id = ?", iamID).Error; err == nil {
			if email := strings.TrimSpace(user.Email); email != "" {
				return email
			}
		}
	}
	if clientID > 0 {
		var cl client.Client
		if err := global.GVA_DB.Select("email").First(&cl, "id = ?", clientID).Error; err == nil {
			return strings.TrimSpace(cl.Email)
		}
	}
	return ""
}

func usedCardHolderEmails(clientID uint) map[string]bool {
	used := map[string]bool{}
	if clientID == 0 || global.GVA_DB == nil {
		return used
	}
	var emails []string
	global.GVA_DB.Model(&modelfinance.CardHolder{}).Where("client_id = ?", clientID).Pluck("email", &emails)
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			used[email] = true
		}
	}
	return used
}

// pickCardHolderEmail 优先返回用户自己的邮箱。已被占用时，在同一邮箱上加随机数字；没有用户邮箱时按姓名生成随机邮箱。
func pickCardHolderEmail(userEmail, firstName, lastName string, used map[string]bool) string {
	if used == nil {
		used = map[string]bool{}
	}
	userEmail = strings.ToLower(strings.TrimSpace(userEmail))
	if emailPattern.MatchString(userEmail) && !used[userEmail] {
		return userEmail
	}
	if emailPattern.MatchString(userEmail) {
		local, domain, _ := strings.Cut(userEmail, "@")
		local = sanitizeEmailLocal(local)
		if local == "" {
			local = "user"
		}
		for i := 0; i < 8; i++ {
			candidate := fmt.Sprintf("%s%d@%s", local, rand.Intn(9000)+1000, domain)
			if !used[candidate] {
				return candidate
			}
		}
	}
	return randomNameEmail(firstName, lastName, used)
}

func randomNameEmail(firstName, lastName string, used map[string]bool) string {
	first := sanitizeEmailLocal(firstName)
	last := sanitizeEmailLocal(lastName)
	if first == "" {
		first = "user"
	}
	if last == "" {
		last = "card"
	}
	for i := 0; i < 8; i++ {
		candidate := fmt.Sprintf("%s%d.%s%d@%s", first, rand.Intn(999)+1, last, rand.Intn(999)+1, randomEmailDomains[rand.Intn(len(randomEmailDomains))])
		if !used[candidate] {
			return candidate
		}
	}
	return fmt.Sprintf("%s%d@gmail.com", first, rand.Intn(900000)+100000)
}

func sanitizeEmailLocal(value string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
