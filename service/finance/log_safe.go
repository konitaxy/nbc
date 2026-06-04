package finance

import (
	"gitlab.com/ucard/logredact"
	"gitlab.com/ucard/model/finance"
	"go.uber.org/zap"
)

// ZapPixielCard 记录卡信息日志时脱敏 CVV。
func ZapPixielCard(card finance.PixielCard) zap.Field {
	card.CVV = logredact.CVV
	return zap.Any("card", card)
}
