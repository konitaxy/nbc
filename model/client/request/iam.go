package request

// IAMLoginReq IAM 用户登录请求
type IAMLoginReq struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱
	Password string `json:"password" binding:"required"`    // 密码
}

// SendIAMCodeReq 发送验证码请求
type SendIAMCodeReq struct {
	Email string `json:"email" binding:"required,email"` // 邮箱
}

// CreateIAMUserReq 创建 IAM 用户请求
type CreateIAMUserReq struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱 (登录名)
	Password string `json:"password" binding:"required"`    // 密码
	Nickname string `json:"nickname"`                       // 昵称
	Roles    []uint `json:"roles"`                          // 角色ID列表
}

// UpdateIAMUserReq 更新 IAM 用户请求
type UpdateIAMUserReq struct {
	ID       uint   `json:"id" binding:"required"` // 用户ID
	Password string `json:"password"`              // 密码 (可选，不传则不修改)
	Nickname string `json:"nickname"`              // 昵称
	Status   *int8  `json:"status"`                // 状态
	Roles    []uint `json:"roles"`                 // 角色ID列表
}

// ListIAMUserReq 获取 IAM 用户列表请求
type ListIAMUserReq struct {
	Page     int    `json:"page" form:"page"`
	ClientID uint   `json:"clientId" form:"clientId"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Email    string `json:"email" form:"email"`   // 邮箱筛选
	Status   int8   `json:"status" form:"status"` // 状态筛选
}

// DeleteIAMUserReq 删除 IAM 用户请求
type DeleteIAMUserReq struct {
	ID uint `json:"id" binding:"required"` // 用户ID
}

// UpdateIAMUserStatusReq 修改 IAM 用户状态请求
type UpdateIAMUserStatusReq struct {
	ID     uint `json:"id" binding:"required"` // 用户ID
	Status int8 `json:"status"`                // 状态 1启用 0禁用
}

// ResetIAMPasswordReq Reset IAM user password request
type ResetIAMPasswordReq struct {
	ID       uint   `json:"id" binding:"required"`       // User ID
	Password string `json:"password" binding:"required"` // New password
}
