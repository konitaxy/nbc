package system

import (
	"reflect"

	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var AuthoritiesMenus = new(authoritiesMenus)

type authoritiesMenus struct{}

func (a *authoritiesMenus) TableName() string {
	var entity AuthorityMenus
	return entity.TableName()
}

func (a *authoritiesMenus) Initialize() error {
	entities := []AuthorityMenus{

		{BaseMenuId: 3, AuthorityId: "888"},
		{BaseMenuId: 31, AuthorityId: "888"},
		{BaseMenuId: 32, AuthorityId: "888"},
		{BaseMenuId: 33, AuthorityId: "888"},
		{BaseMenuId: 34, AuthorityId: "888"},
		{BaseMenuId: 35, AuthorityId: "888"},
		{BaseMenuId: 36, AuthorityId: "888"},
		{BaseMenuId: 37, AuthorityId: "888"},
		{BaseMenuId: 94, AuthorityId: "888"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, a.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (a *authoritiesMenus) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("sys_base_menu_id = ? AND sys_authority_authority_id = ?", 3, "888").First(&AuthorityMenus{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}

type AuthorityMenus struct {
	AuthorityId string `gorm:"column:sys_authority_authority_id"`
	BaseMenuId  uint   `gorm:"column:sys_base_menu_id"`
}

func (a *AuthorityMenus) TableName() string {
	var entity system.SysAuthority
	types := reflect.TypeOf(entity)
	if s, o := types.FieldByName("SysBaseMenus"); o {
		m1 := schema.ParseTagSetting(s.Tag.Get("gorm"), ";")
		return m1["MANY2MANY"]
	}
	return ""
}
