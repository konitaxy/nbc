package finance

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/service/credit_provider/gzy"
	"gitlab.com/ucard/utils"
)

// SandBoxTransaction 光子沙箱交易模拟；requestId 自动生成，CVV/有效期从本地卡记录读取（可覆盖）。
func (f *FinanceService) SandBoxTransaction(req request.SandBoxTransactionSimReq) (*gzy.SandBoxTransactionResponse, error) {
	cardID := strings.TrimSpace(req.CardID)
	if cardID == "" {
		return nil, fmt.Errorf("cardId is required")
	}
	card, err := f.GetCardByCardID(cardID)
	if err != nil || card.ID == 0 {
		return nil, fmt.Errorf("card not found")
	}
	bin, _ := f.getCardBinForPixielCard(&card)
	if bin != nil && strings.TrimSpace(bin.Channel) != string(constant.Channel_Gzy) {
		return nil, fmt.Errorf("sandbox transaction only supports gzy channel cards")
	}

	cvv := strings.TrimSpace(card.CVV)
	if cvv == "" {
		return nil, fmt.Errorf("card cvv is empty, sync card detail first")
	}
	expiry := gzy.ExpiryToSandBoxMMYY(card.Expirey)
	if expiry == "" {
		return nil, fmt.Errorf("card expiration is empty, sync card detail first")
	}

	txnType := strings.ToLower(strings.TrimSpace(req.TxnType))
	if txnType == "" {
		txnType = "auth"
	}
	txnCurrency := strings.TrimSpace(req.TxnCurrency)
	if txnCurrency == "" {
		txnCurrency = string(card.Currency)
	}
	if txnCurrency == "" {
		txnCurrency = "USD"
	}
	amount := req.TxnAmount
	if !amount.IsPositive() {
		amount = decimal.NewFromFloat(1.00)
	}

	if err := gzy.EnsureAccessToken(); err != nil {
		return nil, err
	}

	body := gzy.SandBoxTransactionRequest{
		RequestID:           utils.GenerateID(constant.OrderPrefix_Card_Recharge),
		CardID:              cardID,
		CVV:                 cvv,
		ExpirationDate:      expiry,
		OriginTransactionID: strings.TrimSpace(req.OriginTransactionID),
		TxnCurrency:         txnCurrency,
		TxnAmount:           amount,
		TxnType:             txnType,
		MCC:                 defaultSandBoxStr(req.MCC, "5411"),
		MerchantName:        defaultSandBoxStr(req.MerchantName, "Amazon"),
		MerchantCountry:     defaultSandBoxStr(req.MerchantCountry, "US"),
		MerchantCity:        defaultSandBoxStr(req.MerchantCity, "Newyork"),
		MerchantPostcode:    defaultSandBoxStr(req.MerchantPostcode, "10001"),
	}
	return gzy.NewGzy().SandBoxTransaction(body)
}

func defaultSandBoxStr(v, def string) string {
	if strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return def
}

func (f *FinanceService) getCardBinForPixielCard(card *finance.PixielCard) (*finance.CardBin, error) {
	if card == nil || strings.TrimSpace(card.CardBinID) == "" {
		return nil, fmt.Errorf("card bin id empty")
	}
	var bin finance.CardBin
	if err := global.GVA_DB.Where("card_bin_id = ?", card.CardBinID).First(&bin).Error; err != nil {
		return nil, err
	}
	return &bin, nil
}
