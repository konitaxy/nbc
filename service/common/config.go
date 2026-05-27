package common

import (
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
)

type ConfigService struct {
}

func (ConfigService) Get(key string) (cfg common.AdminConfig, err error) {
	err = global.GVA_DB.Debug().First(&cfg, "kay = ?", key).Error
	return
}

func (ConfigService) Set(cfg *common.AdminConfig) (err error) {
	err = global.GVA_DB.Save(cfg).Error
	return
}

func (ConfigService) Delete(key string) (err error) {
	err = global.GVA_DB.Unscoped().Delete(&common.AdminConfig{}, "kay = ?", key).Error
	return
}
