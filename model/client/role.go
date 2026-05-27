package client

import "gitlab.com/ucard/global"

// ClientRole 角色定义表
type ClientRole struct {
	global.GVA_MODEL
	RoleName    string `gorm:"column:role_name;type:varchar(64);not null" json:"roleName" form:"roleName"`       // 角色名称
	RoleCode    string `gorm:"column:role_code;type:varchar(64);uniqueIndex" json:"roleCode" form:"roleCode"`    // 角色代码
	Description string `gorm:"column:description;type:varchar(255)" json:"description" form:"description"`       // 描述
	Status      int8   `gorm:"column:status;default:1" json:"status" form:"status"`                              // 状态: 1启用 0禁用
	Sort        int    `gorm:"column:sort;default:0" json:"sort" form:"sort"`                                    // 排序
}

func (ClientRole) TableName() string {
	return "client_roles"
}

// ClientUserRole 用户角色关联表
type ClientUserRole struct {
	ClientID uint `gorm:"column:client_id;primaryKey;index" json:"clientId" form:"clientId"` // 用户ID
	RoleID   uint `gorm:"column:role_id;primaryKey;index" json:"roleId" form:"roleId"`       // 角色ID
}

func (ClientUserRole) TableName() string {
	return "client_user_roles"
}

// ClientRolePermission 角色权限关联表
type ClientRolePermission struct {
	RoleID       uint `gorm:"column:role_id;primaryKey;index" json:"roleId" form:"roleId"`             // 角色ID
	PermissionID uint `gorm:"column:permission_id;primaryKey;index" json:"permissionId" form:"permissionId"` // 权限ID
}

func (ClientRolePermission) TableName() string {
	return "client_role_permissions"
}

// 默认角色
var DefaultClientRoles = []ClientRole{
	{RoleName: "管理员", RoleCode: "admin", Description: "拥有所有权限", Sort: 1},
	{RoleName: "操作员", RoleCode: "operator", Description: "可以进行卡片操作", Sort: 2},
	{RoleName: "财务", RoleCode: "finance", Description: "可以查看财务信息", Sort: 3},
	{RoleName: "只读", RoleCode: "readonly", Description: "只能查看信息", Sort: 4},
}



