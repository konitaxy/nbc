package client

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gorm.io/gorm"
)

var IAM = new(iam)

type iam struct{}

func (i *iam) TableName() string {
	return "iam_permissions"
}

// 默认 IAM 权限列表
var DefaultIAMPermissions = []client.IAMPermission{
	// ===== 基础功能权限 =====
	{Code: "api:wallet:view", Name: "钱包查看", Type: "api", Action: "wallet:view", Sort: 1},
	{Code: "api:profile:view", Name: "我的信息", Type: "api", Action: "profile:view", Sort: 2},
	{Code: "api:2fa:manage", Name: "Google认证", Type: "api", Action: "2fa:manage", Sort: 3},
	{Code: "api:verify:setting", Name: "验证码设置", Type: "api", Action: "verify:setting", Sort: 4},
	{Code: "api:notice:view", Name: "查看公告", Type: "api", Action: "notice:view", Sort: 5},

	// ===== 卡片查看权限 =====
	{Code: "api:card:view", Name: "查看卡信息(不含CVV)", Type: "api", Action: "card:view", Sort: 10},
	{Code: "api:card:log", Name: "开卡日志", Type: "api", Action: "card:log", Sort: 11},
	{Code: "api:transaction:view", Name: "查看交易明细", Type: "api", Action: "transaction:view", Sort: 12},
	{Code: "api:cancel_record:view", Name: "查看销卡信息", Type: "api", Action: "cancel_record:view", Sort: 13},
	{Code: "api:recharge_record:view", Name: "查看充值/退款信息", Type: "api", Action: "recharge_record:view", Sort: 14},

	// ===== 卡片管理权限 =====
	{Code: "api:card:create", Name: "常规开卡", Type: "api", Action: "card:create", Sort: 20},
	{Code: "api:card:quick_create", Name: "快速开卡", Type: "api", Action: "card:quick_create", Sort: 21},
	{Code: "api:card:cancel", Name: "批量销卡", Type: "api", Action: "card:cancel", Sort: 22},
	{Code: "api:card:batch_cancel", Name: "批量销卡", Type: "api", Action: "card:batch_cancel", Sort: 23},
	{Code: "api:card:freeze", Name: "冻结解冻", Type: "api", Action: "card:freeze", Sort: 24},
	{Code: "api:card:batch_freeze", Name: "批量冻结解冻", Type: "api", Action: "card:batch_freeze", Sort: 25},
	{Code: "api:card:limit", Name: "设置限额", Type: "api", Action: "card:limit", Sort: 26},
	{Code: "api:card:remark", Name: "编辑卡备注", Type: "api", Action: "card:remark", Sort: 27},
	{Code: "api:card:update", Name: "更新卡信息", Type: "api", Action: "card:update", Sort: 28},
	{Code: "api:cardholder:view", Name: "查看持卡人信息", Type: "api", Action: "cardholder:view", Sort: 29},
	{Code: "api:cardholder:update", Name: "修改持卡人信息", Type: "api", Action: "cardholder:update", Sort: 30},
	{Code: "api:cardholder:disable", Name: "停用持卡人", Type: "api", Action: "cardholder:disable", Sort: 31},

	// ===== 敏感操作权限 =====
	{Code: "api:card:cvv", Name: "查看卡CVV", Type: "api", Action: "card:cvv", Sort: 40},
	{Code: "api:otp:view", Name: "查看OTP", Type: "api", Action: "otp:view", Sort: 41},

	// ===== 财务操作权限 =====
	{Code: "api:card:recharge", Name: "卡充值", Type: "api", Action: "card:recharge", Sort: 50},
	{Code: "api:card:refund", Name: "卡金额退回", Type: "api", Action: "card:refund", Sort: 51},
	{Code: "api:fee:view", Name: "查看手续费明细", Type: "api", Action: "fee:view", Sort: 52},
	{Code: "api:commission:view", Name: "查看分成信息", Type: "api", Action: "commission:view", Sort: 53},

	// ===== 数据导出权限 =====
	{Code: "api:card:export", Name: "导出卡信息", Type: "api", Action: "card:export", Sort: 60},
	{Code: "api:transaction:export", Name: "导出交易明细", Type: "api", Action: "transaction:export", Sort: 61},
	{Code: "api:fee:export", Name: "导出手续费明细", Type: "api", Action: "fee:export", Sort: 62},
	{Code: "api:cancel_record:export", Name: "导出销卡信息", Type: "api", Action: "cancel_record:export", Sort: 63},
	{Code: "api:recharge_record:export", Name: "导出充值/退款信息", Type: "api", Action: "recharge_record:export", Sort: 64},
	{Code: "api:commission:export", Name: "导出分成信息", Type: "api", Action: "commission:export", Sort: 65},
}

// Default IAM Roles (6 roles)
var DefaultIAMRoles = []client.IAMRole{
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 1},
		RoleCode:    "role_base",
		RoleName:    "Basic_Features",
		Description: "Wallet_view,_profile,_Google_authentication,_verification_code_settings,_view_announcements",
		IsDefault:   true,
		Sort:        1,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 2},
		RoleCode:    "role_card_view",
		RoleName:    "Card_View",
		Description: "View_card_list,_card_details_(without_CVV),_transaction_details,_cancellation_records,_recharge_records",
		IsDefault:   false,
		Sort:        2,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 3},
		RoleCode:    "role_card_manage",
		RoleName:    "Card_Management",
		Description: "Create_card,_cancel_card,_freeze/unfreeze,_set_limit,_edit_notes,_cardholder_management",
		IsDefault:   false,
		Sort:        3,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 4},
		RoleCode:    "role_card_sensitive",
		RoleName:    "Sensitive_Operations",
		Description: "View_CVV,_view_OTP",
		IsDefault:   false,
		Sort:        4,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 5},
		RoleCode:    "role_finance",
		RoleName:    "Financial_Operations",
		Description: "Card_recharge,_card_refund,_view_fees,_view_commission_info",
		IsDefault:   false,
		Sort:        5,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 6},
		RoleCode:    "role_export",
		RoleName:    "Data_Export",
		Description: "Export_card_info,_transaction_details,_fees,_cancellation_records,_recharge_records,_commission_info",
		IsDefault:   false,
		Sort:        6,
	},
	{
		GVA_MODEL:   global.GVA_MODEL{ID: 7},
		RoleCode:    "role_wallet",
		RoleName:    "Wallet_Operations",
		Description: "Wallet_recharge,_wallet_withdraw",
		IsDefault:   false,
		Sort:        7,
	},
}

// 角色权限映射
var RolePermissionMap = map[string][]string{
	// 基础功能 - 所有子账号默认拥有
	"role_base": {
		"api:wallet:view",
		"api:profile:view",
		"api:2fa:manage",
		"api:verify:setting",
		"api:notice:view",
	},
	// 卡片查看 - 只读权限
	"role_card_view": {
		"api:card:view",
		"api:card:log",
		"api:transaction:view",
		"api:cancel_record:view",
		"api:recharge_record:view",
	},
	// 卡片管理 - 操作权限
	"role_card_manage": {
		"api:card:create",
		"api:card:quick_create",
		"api:card:cancel",
		"api:card:batch_cancel",
		"api:card:freeze",
		"api:card:batch_freeze",
		"api:card:limit",
		"api:card:remark",
		"api:card:update",
		"api:cardholder:view",
		"api:cardholder:update",
		"api:cardholder:disable",
	},
	// 敏感操作 - CVV、OTP等
	"role_card_sensitive": {
		"api:card:cvv",
		"api:otp:view",
	},
	// 财务操作 - 充值、退款等
	"role_finance": {
		"api:card:recharge",
		"api:card:refund",
		"api:fee:view",
		"api:commission:view",
	},
	// 数据导出 - 所有导出功能
	"role_export": {
		"api:card:export",
		"api:transaction:export",
		"api:fee:export",
		"api:cancel_record:export",
		"api:recharge_record:export",
		"api:commission:export",
	},
	// 钱包操作 - 充值、提现
	"role_wallet": {
		"api:wallet:recharge",
		"api:wallet:withdraw",
	},
}

// Initialize 初始化 IAM 默认数据
func (i *iam) Initialize() error {
	// 初始化权限
	if err := i.initPermissions(); err != nil {
		return err
	}
	// 初始化角色
	if err := i.initRoles(); err != nil {
		return err
	}
	return nil
}

// initPermissions 初始化默认权限
func (i *iam) initPermissions() error {
	for _, p := range DefaultIAMPermissions {
		var existing client.IAMPermission
		if err := global.GVA_DB.Where("code = ?", p.Code).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := global.GVA_DB.Create(&p).Error; err != nil {
					return errors.Wrapf(err, "创建权限 %s 失败", p.Code)
				}
			} else {
				return err
			}
		}
	}
	return nil
}

// initRoles 初始化默认角色
func (i *iam) initRoles() error {
	for _, r := range DefaultIAMRoles {
		var existing client.IAMRole
		if err := global.GVA_DB.Where("role_code = ? AND client_id = 0", r.RoleCode).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 根据角色代码获取对应的权限
				if permCodes, ok := RolePermissionMap[r.RoleCode]; ok {
					var permissions []client.IAMPermission
					global.GVA_DB.Where("code IN ?", permCodes).Find(&permissions)
					r.Permissions = permissions
				}
				if err := global.GVA_DB.Create(&r).Error; err != nil {
					return errors.Wrapf(err, "创建角色 %s 失败", r.RoleCode)
				}
			} else {
				return err
			}
		}
	}
	return nil
}

// CheckDataExist 检查数据是否存在
func (i *iam) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("code = ?", "api:wallet:view").First(&client.IAMPermission{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
