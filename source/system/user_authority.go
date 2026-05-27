package system

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gorm.io/gorm"
)

var UserAuthority = new(userAuthority)

type userAuthority struct{}

func (a *userAuthority) TableName() string {
	var entity system.SysUseAuthority
	return entity.TableName()
}

func (a *userAuthority) Initialize() error {
	entities := []system.SysUseAuthority{
		{SysUserId: 3, SysAuthorityAuthorityId: "888"},
		{SysUserId: 31, SysAuthorityAuthorityId: "888"},
		{SysUserId: 32, SysAuthorityAuthorityId: "888"},
		{SysUserId: 33, SysAuthorityAuthorityId: "888"},
		{SysUserId: 34, SysAuthorityAuthorityId: "888"},
		{SysUserId: 35, SysAuthorityAuthorityId: "888"},
		{SysUserId: 36, SysAuthorityAuthorityId: "888"},
		{SysUserId: 37, SysAuthorityAuthorityId: "888"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, a.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (a *userAuthority) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", 3, "888").First(&system.SysUseAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
