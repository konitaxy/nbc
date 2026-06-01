package system

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gorm.io/gorm"
)

var BaseMenu = new(menu)

type menu struct{}

func (m *menu) TableName() string {
	return "sys_base_menus"
}

func (m *menu) Initialize() error {
	entities := []system.SysBaseMenu{
		{GVA_MODEL: global.GVA_MODEL{ID: 3}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 99, Meta: system.Meta{Title: "超级管理员", Icon: "user"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 31}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: system.Meta{Title: "角色管理", Icon: "avatar"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 32}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: system.Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 33}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: system.Meta{Title: "api管理", Icon: "platform", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 34}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: system.Meta{Title: "用户管理", Icon: "coordinate"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 35}, MenuLevel: 0, Hidden: true, ParentId: "0", Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 4, Meta: system.Meta{Title: "个人信息", Icon: "message"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 17}, MenuLevel: 0, Hidden: false, ParentId: "14", Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 3, Meta: system.Meta{Title: "系统配置", Icon: "operation"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 36}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: system.Meta{Title: "操作历史", Icon: "pie-chart"}},

		//添加新页面
		{GVA_MODEL: global.GVA_MODEL{ID: 4}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dashbard", Name: "dashboard", Component: "view/admin/dashboard/index.vue", Sort: 1, Meta: system.Meta{Title: "Dashboard", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 41}, MenuLevel: 0, Hidden: false, ParentId: "4", Path: "summary", Name: "summary", Component: "view/admin/dashboard/index.vue", Sort: 1, Meta: system.Meta{Title: "Summary", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 42}, MenuLevel: 0, Hidden: false, ParentId: "4", Path: "anylisis", Name: "transactionAnalysis", Component: "view/admin/dashboard/list/index.vue", Sort: 2, Meta: system.Meta{Title: "Customer Analysis", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 5}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "kyc", Name: "kyc", Component: "view/admin/kyc/index.vue", Sort: 2, Meta: system.Meta{Title: "KYC", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 51}, MenuLevel: 0, Hidden: false, ParentId: "5", Path: "dueDiligence", Name: "dueDiligence", Component: "view/admin/kyc/dueDiligence/index.vue", Sort: 1, Meta: system.Meta{Title: "Due Diligence", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 6}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "customer", Name: "customer", Component: "view/admin/customer/index.vue", Sort: 3, Meta: system.Meta{Title: "Customer", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 61}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "list", Name: "customerList", Component: "view/admin/customer/list/index.vue", Sort: 1, Meta: system.Meta{Title: "Customer List", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 62}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "report", Name: "monthReport", Component: "view/admin/customer/report/index.vue", Sort: 2, Meta: system.Meta{Title: "Month Report", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 7}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "transactions", Name: "transactions", Component: "view/admin/transactions/index.vue", Sort: 4, Meta: system.Meta{Title: "Transactions", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 71}, MenuLevel: 0, Hidden: false, ParentId: "7", Path: "list", Name: "transactionsList", Component: "view/admin/transactions/list/index.vue", Sort: 1, Meta: system.Meta{Title: "Transactions List", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 72}, MenuLevel: 0, Hidden: false, ParentId: "7", Path: "negative", Name: "negativeBalance", Component: "view/admin/transactions/negativeBalance/index.vue", Sort: 2, Meta: system.Meta{Title: "Negative Balance", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 73}, MenuLevel: 0, Hidden: false, ParentId: "7", Path: "withdraw", Name: "walletWithdraw", Component: "view/admin/transactions/walletWithdraw/index.vue", Sort: 3, Meta: system.Meta{Title: "Wallet Withdraw", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 8}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "cardCenter", Name: "cardCenter", Component: "view/admin/cardCenter/index.vue", Sort: 5, Meta: system.Meta{Title: "Card Center", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 81}, MenuLevel: 0, Hidden: false, ParentId: "8", Path: "cardBin", Name: "cardBin", Component: "view/admin/cardCenter/cardBin/index.vue", Sort: 1, Meta: system.Meta{Title: "Card Bin", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 82}, MenuLevel: 0, Hidden: false, ParentId: "8", Path: "management", Name: "cardManagement", Component: "view/admin/cardCenter/cardManagement/index.vue", Sort: 2, Meta: system.Meta{Title: "Card Management", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 83}, MenuLevel: 0, Hidden: false, ParentId: "8", Path: "holder", Name: "cardHolder", Component: "view/admin/cardCenter/cardHolder/index.vue", Sort: 3, Meta: system.Meta{Title: "Card Holder", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 9}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "configue", Name: "configue", Component: "view/admin/configue/index.vue", Sort: 6, Meta: system.Meta{Title: "Configue", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 91}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "cardFee", Name: "cardFee", Component: "view/admin/configue/cardFee/index.vue", Sort: 1, Meta: system.Meta{Title: "Card Fee", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 92}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "inboundFee", Name: "inboundFee", Component: "view/admin/configue/inboundFee/index.vue", Sort: 2, Meta: system.Meta{Title: "Inbound Fee", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 93}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "setting", Name: "inboundAndDepositSetting", Component: "view/admin/configue/inboundSetting/index.vue", Sort: 3, Meta: system.Meta{Title: "Inbound Settings", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 94}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "chain", Name: "chain", Component: "view/admin/configue/chain/index.vue", Sort: 4, Meta: system.Meta{Title: "Chain", Icon: "pie-chart"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 10}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "snsCode", Name: "smsCode", Component: "view/admin/smsCode/index.vue", Sort: 7, Meta: system.Meta{Title: "SMS Code", Icon: "pie-chart"}},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil { // 创建 model.User 初始化数据
		return errors.Wrap(err, m.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (m *menu) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("path = ?", "admin").First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
