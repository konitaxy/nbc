package system

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gorm.io/gorm"
)

var Authority = new(authority)

type authority struct{}

func (a *authority) TableName() string {
	return "sys_authorities"
}

func (a *authority) Initialize() error {
	entities := []system.SysAuthority{
		{AuthorityId: "888", AuthorityName: "普通用户", ParentId: "0", DefaultRouter: "artist"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrapf(err, "%s表数据初始化失败!", a.TableName())
	}
	return nil
}

func (a *authority) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("authority_id = ?", "888").First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
