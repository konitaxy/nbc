package utils

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ucard/global"
	systemReq "gitlab.com/ucard/model/system/request"
)

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := c.Request.Header.Get("x-token")
	if token == "" {
		return nil, nil
	}
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.ID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.ID
	}
}

// 从Gin的Context中获取从jwt解析出来的用户ID
func GetRoles(c *gin.Context) []uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return []uint{}
		} else {
			return cl.Roles
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Roles
	}
}

func GetTenantID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			if cl.TenantID == 0 {
				return cl.ID
			}
			return cl.TenantID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		if waitUse.TenantID == 0 {
			return waitUse.ID
		}
		return waitUse.TenantID
	}
}

// 从Gin的Context中获取从jwt解析出来的用户ID
// iamID,tenantID,isIAM
func GetUserAndTenantID(c *gin.Context) (uint, uint, bool) {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0, 0, false
		} else {
			if cl.IsIAM {
				return cl.ID, cl.TenantID, true
			}
			return cl.ID, cl.ID, false
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		if waitUse.IsIAM {
			return waitUse.ID, waitUse.TenantID, true
		}
		return waitUse.ID, waitUse.ID, false
	}
}

func GetUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Username
	}
}

func GetUserEmail(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Email
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Email
	}
}

// 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.UUID
	}
}

// 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.AuthorityId
	}
}

// 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

// IsIAM returns true if the current user is an IAM user
func IsIAM(c *gin.Context) bool {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return false
		} else {
			return cl.IsIAM
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.IsIAM
	}
}

// GetIAMID returns the IAM user ID from context
func GetIAMID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.ID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.ID
	}
}
