package utils

import (
	"fmt"
	"strings"

	"gitlab.com/ucard/global"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const teml = `<!DOCTYPE html>
<html>
<head>
  <title>Verification Code Email</title>
  <style type="text/css">
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f4;
      margin: 0;
      padding: 0;
    }
    .email-container {
      width: 100%;
      max-width: 600px;
      margin: auto;
      background-color: #ffffff;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 20px;
    }
    .header img {
      max-width: 150px;
    }
    .content {
      margin-bottom: 20px;
    }
    .code {
      font-size: 24px;
      font-weight: bold;
      color: #333;
      text-align: center;
      margin-top: 20px;
    }
    .footer {
      color: #777;
      font-size: 12px;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="email-container">
    <!-- Header -->
    <div class="header">
      <img src="https://metalposterpro.s3.us-east-1.amazonaws.com/static/artist.png" alt="Company Logo">
      <h2>{{.CompanyName}}</h2>
    </div>

    <!-- Content -->
    <div class="content">
      <p>Dear User,</p>
      <p>Thank you! To verify your email address, please use the verification code below.</p>
    </div>

    <!-- Verification Code -->
    <div class="code">
      Your verification code is: <span style="color:#ff6600;">{{.Code}}</span>
    </div>

    <!-- Footer -->
    <div class="footer">
      <p>If you did not request this verification code, please ignore this email.</p>
      <p>© {{.Year}} {{.CompanyName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>`

// func SendEmail(to []string, subject, body string) error {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", "melong_tech@qq.com")
// 	m.SetHeader("To", to...)
// 	m.SetHeader("Subject", subject)
// 	m.SetHeader("Cc", "melong_tech@qq.com")

// 	var bs bytes.Buffer
// 	tmpData := struct {
// 		Code        string
// 		CompanyName string
// 		Year        uint
// 	}{Code: body, CompanyName: "Melong Tech", Year: uint(time.Now().Year())}
// 	t := template.Must(template.New("email").Parse(teml))
// 	t.Execute(&bs, tmpData)
// 	m.SetBody("text/html", bs.String())

// 	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}
// 	return nil
// }

func SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", "Newbeecard.com", global.GVA_CONFIG.Email.From))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	// m.SetHeader("Cc", global.GVA_CONFIG.Email.From)
	m.SetHeader("X-MC-Template", "verify_code")
	m.SetHeader("X-MC-MergeVars", fmt.Sprintf(`{ "Code": "%s"}`, body))
	m.SetHeader("X-MC-MergeLanguage", "handlebars")

	// m.SetBody("text/plain", fmt.Sprintf("<html><body>%s</body></html>", body))

	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func SendApplyReceivedEmail(to []string, subject, costDay, ApplyName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", "Metalposter.com", global.GVA_CONFIG.Email.From))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetHeader("Reply-To", global.GVA_CONFIG.Email.ReplyTo)
	// m.SetHeader("Cc", global.GVA_CONFIG.Email.From)
	m.SetHeader("X-MC-Template", "ApplicationReceived")
	m.SetHeader("X-MC-MergeVars", fmt.Sprintf(`{ "CostDay": "%s","ApplyName":"%s","PlatformName":"%s","SupportMail":"%s"}`, costDay, ApplyName, "Metalposter", "support@metalposter.com"))
	m.SetHeader("X-MC-MergeLanguage", "handlebars")

	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func SendApplyAcceptEmail(to []string, subject, InviteLink, ApplyName string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", "Metalposter.com", global.GVA_CONFIG.Email.From))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetHeader("Reply-To", global.GVA_CONFIG.Email.ReplyTo)
	// m.SetHeader("Cc", global.GVA_CONFIG.Email.From)
	m.SetHeader("X-MC-Template", "ApplicationApproved")
	m.SetHeader("X-MC-MergeVars", fmt.Sprintf(`{ "InviteLink":"%s","ApplyName":"%s","PlatformName":"%s","SupportMail":"%s"}`, InviteLink, ApplyName, "Metalposter", "support@metalposter.com"))
	m.SetHeader("X-MC-MergeLanguage", "handlebars")

	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
	if err := d.DialAndSend(m); err != nil {
		global.GVA_LOG.Error("SendApplyAcceptEmail Error", zap.Error(err))
		return err
	}
	global.GVA_LOG.Sugar().Infof("SendApplyAcceptEmail Success[ %s %s %s ]", to[0], InviteLink, ApplyName)
	return nil
}

// SendNotifyEmail 向配置 system.admin 发送 HTML 通知邮件（支持多个邮箱；不走 Mandrill 模板）。
func SendNotifyEmail(subject, htmlBody string) error {
	to := resolveAdminNotifyEmails()
	if len(to) == 0 {
		return fmt.Errorf("system.admin is empty")
	}

	fromName := strings.TrimSpace(global.GVA_CONFIG.Email.Nickname)
	if fromName == "" {
		fromName = "Newbeecard.com"
	}
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, global.GVA_CONFIG.Email.From))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	if reply := strings.TrimSpace(global.GVA_CONFIG.Email.ReplyTo); reply != "" {
		m.SetHeader("Reply-To", reply)
	}
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(global.GVA_CONFIG.Email.Host, global.GVA_CONFIG.Email.Port, global.GVA_CONFIG.Email.From, global.GVA_CONFIG.Email.Secret)
	if err := d.DialAndSend(m); err != nil {
		global.GVA_LOG.Error("SendNotifyEmail failed", zap.Error(err), zap.Strings("to", to), zap.String("subject", subject))
		return err
	}
	return nil
}

func resolveAdminNotifyEmails() []string {
	seen := make(map[string]struct{})
	var to []string
	for _, part := range global.GVA_CONFIG.System.Admin {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		to = append(to, s)
	}
	return to
}
