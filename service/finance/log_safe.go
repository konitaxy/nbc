package finance

import (
	"gitlab.com/ucard/logredact"
	"gitlab.com/ucard/model/finance"
	"go.uber.org/zap"
)

// ZapPixielCard 记录卡信息日志时脱敏 CVV、有效期等敏感字段。
func ZapPixielCard(card finance.PixielCard) zap.Field {
	card.CVV = logredact.Redacted
	card.Expirey = logredact.Redacted
	return zap.Any("card", card)
}
