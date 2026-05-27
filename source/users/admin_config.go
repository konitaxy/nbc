package users

import (
	"github.com/pkg/errors"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/utils/model"
	"gorm.io/gorm"
)

var AdminConfig = new(adminConfig)

type adminConfig struct{}

func (a *adminConfig) TableName() string {
	return "admin_configs"
}

func (a *adminConfig) Initialize() error {
	entities := []common.AdminConfig{
		{GVA_MODEL: global.GVA_MODEL{ID: 1}, Kay: "tronAddress", StringValue: model.NewString("lakjsdlkfjalkdjf"), ValueType: "string"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, a.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (a *adminConfig) CheckDataExist() bool {

	return !errors.Is(global.GVA_DB.Where("id = ?", 1).First(&common.AdminConfig{}).Error, gorm.ErrRecordNotFound)
}
