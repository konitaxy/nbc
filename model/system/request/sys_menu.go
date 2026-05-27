package request

import (
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
)

// Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []system.SysBaseMenu `json:"menus"`
	AuthorityId string               `json:"authorityId"` // 角色ID
}

func DefaultMenu() []system.SysBaseMenu {
	return []system.SysBaseMenu{{
		GVA_MODEL: global.GVA_MODEL{ID: 1},
		ParentId:  "0",
		Path:      "artwork",
		Name:      "artwork",
		Component: "view/artwork/index.vue",
		Sort:      1,
		Meta: system.Meta{
			Title: "作品管理",
			Icon:  "setting",
		},
	}}
}
