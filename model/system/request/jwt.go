package request

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

// roleType : 1 :系统用户 2: 平台用户
type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	TenantID    uint // 租户ID (主账号ID)
	IsIAM       bool // 是否IAM用户
	RoleType    int
	Email       string
	Username    string
	NickName    string
	AuthorityId string
	Admin       *uint //管理员
	IsFreeze    bool
	Roles       []uint
}
