package common

import "gitlab.com/ucard/global"

type OpType uint

const (
	OpType_CardBin_Create OpType = iota + 1
	OpType_CardBin_Dpdate
	OpType_CardBin_Block

	OpType_Fee_User_Config_Set
	OpType_Fee_Global_Config_Set

	OpType_System_Config_Set

	OpType_Wallet_Withdraw_Review
	OpType_Wallet_Recharge_Review
	OpType_Wallet_Recharge_Edit

	OpType_Card_Frozen    // 卡冻结
	OpType_Card_UnFrozen  // 卡解冻
)

type OpLog struct {
	global.GVA_MODEL
	Who    uint   `gorm:"index" json:"who"`
	Name   string `gorm:"index" json:"name"`
	OpType OpType `gorm:"index" json:"opType"`
	Detail string `gorm:"type:text" json:"detail"`
	ObjId  uint   `gorm:"index" json:"objId"`
	Source uint   `gorm:"default:1" json:"source"` //1 后台用户 2 前台用户

}

func (OpLog) TableName() string {
	return "operation_log"
}
