package users

import (
	"gitlab.com/ucard/global"
)

var User = new(user)

type user struct{}

func (a *user) TableName() string {
	return "users"
}

func (a *user) Initialize() error {
	global.GVA_DB.Exec("ALTER TABLE clients AUTO_INCREMENT = 1000")
	return nil
}

func (a *user) CheckDataExist() bool {
	// if errors.Is(global.GVA_DB.Where("id = ?", 1000).First(&profile.User{}).Error, gorm.ErrRecordNotFound) {
	// 	return false
	// }
	return true
}
