package client

import "gitlab.com/ucard/global"

// ClientPermission 权限定义表
type ClientPermission struct {
	global.GVA_MODEL
	Code        string `gorm:"column:code;type:varchar(64);uniqueIndex" json:"code" form:"code"`           // 权限代码
	Name        string `gorm:"column:name;type:varchar(64)" json:"name" form:"name"`                       // 权限名称
	Type        string `gorm:"column:type;type:varchar(16);index" json:"type" form:"type"`                 // 类型: page(页面) / api(接口)
	Path        string `gorm:"column:path;type:varchar(128)" json:"path" form:"path"`                      // 路径: 页面路由或API路径
	Method      string `gorm:"column:method;type:varchar(16)" json:"method" form:"method"`                 // HTTP方法(仅API类型): GET/POST/PUT/DELETE
	ParentID    uint   `gorm:"column:parent_id;index;default:0" json:"parentId" form:"parentId"`           // 父级ID(用于页面菜单层级)
	Sort        int    `gorm:"column:sort;default:0" json:"sort" form:"sort"`                              // 排序
	Description string `gorm:"column:description;type:varchar(255)" json:"description" form:"description"` // 描述

	Children []ClientPermission `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (ClientPermission) TableName() string {
	return "client_permissions"
}

// ClientAccountPermission 账号权限关联表
type ClientAccountPermission struct {
	ClientID     uint `gorm:"column:client_id;primaryKey" json:"clientId"`
	PermissionID uint `gorm:"column:permission_id;primaryKey" json:"permissionId"`
}

func (ClientAccountPermission) TableName() string {
	return "client_account_permissions"
}

// 权限类型常量
const (
	PermissionTypePage = "page" // 页面权限
	PermissionTypeAPI  = "api"  // 接口权限
)

// 默认权限列表
var DefaultClientPermissions = []ClientPermission{
	// 页面权限
	{Code: "page_dashboard", Name: "控制台", Type: PermissionTypePage, Path: "/dashboard", Sort: 1},
	{Code: "page_card_list", Name: "卡片列表", Type: PermissionTypePage, Path: "/card/list", Sort: 2},
	{Code: "page_card_detail", Name: "卡片详情", Type: PermissionTypePage, Path: "/card/detail", Sort: 3},
	{Code: "page_wallet", Name: "钱包管理", Type: PermissionTypePage, Path: "/wallet", Sort: 4},
	{Code: "page_transaction", Name: "交易记录", Type: PermissionTypePage, Path: "/transaction", Sort: 5},

	// 接口权限 - 卡片相关
	{Code: "api_card_list", Name: "查看卡列表", Type: PermissionTypeAPI, Path: "/card/list", Method: "GET", Sort: 10},
	{Code: "api_card_detail", Name: "查看卡详情", Type: PermissionTypeAPI, Path: "/card/detail", Method: "GET", Sort: 11},
	{Code: "api_card_add", Name: "开卡", Type: PermissionTypeAPI, Path: "/card/add", Method: "POST", Sort: 12},
	{Code: "api_card_recharge", Name: "卡充值", Type: PermissionTypeAPI, Path: "/card/recharge", Method: "POST", Sort: 13},
	{Code: "api_card_withdraw", Name: "卡提现", Type: PermissionTypeAPI, Path: "/card/withdraw", Method: "POST", Sort: 14},
	{Code: "api_card_cancel", Name: "批量销卡", Type: PermissionTypeAPI, Path: "/card/cancel", Method: "POST", Sort: 15},

	// 接口权限 - 钱包相关
	{Code: "api_wallet_info", Name: "查看钱包", Type: PermissionTypeAPI, Path: "/wallet/info", Method: "GET", Sort: 20},
	{Code: "api_wallet_withdraw", Name: "钱包提现", Type: PermissionTypeAPI, Path: "/wallet/withdraw/apply", Method: "POST", Sort: 21},
	{Code: "api_wallet_history", Name: "钱包记录", Type: PermissionTypeAPI, Path: "/wallet/history", Method: "GET", Sort: 22},
}
