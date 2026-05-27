package response

import (
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/system"
)

type ClientRes struct {
	AuthorityId string                `json:"authorityId"`
	Authorities []system.SysAuthority `json:"authorities"`
	Authority   system.SysAuthority   `json:"authority"`
	client.Client
}
type LoginRes struct {
	Token     string    `json:"token"`
	Client    ClientRes `json:"user"`
	ExpiresAt int64     `json:"expiresAt"`
}

func ToUserRes(er client.Client) ClientRes {

	// if er.CouponCodeList == nil {
	// 	er.CouponCodeList = make([]user.CouponCode, 0)
	// }
	if er.Wallet == nil {
		er.Wallet = &client.Wallet{
			Balance:  decimal.Zero,
			Currency: constant.USD,
		}
	}
	if er.VerifySetting == nil {
		er.VerifySetting = make(common.MapStringBool)
		for _, v := range client.Default_Verify_Setting {
			er.VerifySetting[v.Key] = v.Value
		}
	}

	return ClientRes{
		AuthorityId: "1618",
		Authorities: []system.SysAuthority{},
		Authority: system.SysAuthority{
			AuthorityId:   "1618",
			AuthorityName: "artist",
			DefaultRouter: "myWallet",
			ParentId:      "0",
		},
		Client: er,
	}
}

// IAMUserRes IAM用户响应结构
type IAMUserRes struct {
	AuthorityId string                `json:"authorityId"`
	Authorities []system.SysAuthority `json:"authorities"`
	Authority   system.SysAuthority   `json:"authority"`
	client.IAMUser
}

// ToIAMUserRes 将 IAMUser 转换为响应结构
func ToIAMUserRes(user client.IAMUser) IAMUserRes {
	if user.Wallet == nil {
		user.Wallet = &client.Wallet{
			Balance:  decimal.Zero,
			Currency: constant.USD,
		}
	}
	if user.VerifySetting == nil {
		user.VerifySetting = make(common.MapStringBool)
		for _, v := range client.Default_IAM_Verify_Setting {
			user.VerifySetting[v.Key] = v.Value
		}
	}

	return IAMUserRes{
		AuthorityId: "1618",
		Authorities: []system.SysAuthority{},
		Authority: system.SysAuthority{
			AuthorityId:   "1618",
			AuthorityName: "iam_user",
			DefaultRouter: "myWallet",
			ParentId:      "0",
		},
		IAMUser: user,
	}
}
