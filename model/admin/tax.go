package admin

import "gitlab.com/ucard/global"

type CountryTax struct {
	global.GVA_MODEL
	Country string `gorm:"column:country;not null;index" json:"country" form:"country"`
	VAT     uint64 `gorm:"column:vat;default:0" json:"vat" form:"vat"`
	Tariff  uint64 `gorm:"column:tariff;default:0" json:"tariff" form:"tariff"`
	Code    string `gorm:"column:code;not null;index" json:"code" form:"code"`
}

func (CountryTax) TableName() string {
	return "country_tax"
}
