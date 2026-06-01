package system

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gorm.io/gorm"
)

// ChainMenuAuthority 为角色 888 补充 Chain 菜单（ID 94）权限；已有库启动时自动执行。
var ChainMenuAuthority = new(chainMenuAuthority)

type chainMenuAuthority struct{}

func (c *chainMenuAuthority) TableName() string {
	return AuthoritiesMenus.TableName()
}

func (c *chainMenuAuthority) Initialize() error {
	return global.GVA_DB.Create(&AuthorityMenus{BaseMenuId: 94, AuthorityId: "888"}).Error
}

func (c *chainMenuAuthority) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("sys_base_menu_id = ? AND sys_authority_authority_id = ?", 94, "888").First(&AuthorityMenus{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
