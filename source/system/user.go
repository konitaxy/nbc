package system

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
	"gorm.io/gorm"
)

var User = new(user)

type user struct{}

func (u *user) TableName() string {
	return "sys_users"
}

func (u *user) Initialize() error {
	entities := []system.SysUser{
		{UUID: uuid.NewV4(), Username: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "超级管理员", HeaderImg: "https://i.pinimg.com/originals/09/97/dd/0997ddb8019c30e5c89eef4fc39f11fd.jpg", AuthorityId: "888", Phone: "18060479363", Email: "jojo@melong.sg", InviteCode: "888888"},
		// {UUID: uuid.NewV4(), Username: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "超级管理员", HeaderImg: "https://i.pinimg.com/originals/09/97/dd/0997ddb8019c30e5c89eef4fc39f11fd.jpg", AuthorityId: "888", Phone: "18060479363", Email: "jojo@melong.sg"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, u.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (u *user) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("username = ?", "admin").First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
