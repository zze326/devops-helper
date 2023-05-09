package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type Menu struct {
	gormc.Model
	Name     string         `gorm:"comment:名称" json:"name"`
	Icon     string         `gorm:"comment:图标" json:"icon"`
	RouteID  int            `gorm:"comment:路由ID" json:"route_id"`
	Route    *FrontendRoute `gorm:"foreignKey:RouteID" json:"route"`
	ParentID int            `gorm:"comment:父级菜单ID" json:"parent_id"`
	Children []*Menu        `gorm:"foreignKey:ParentID" json:"children"`
	Sort     int            `gorm:"comment:排序" json:"sort"`
	Enabled  bool           `gorm:"comment:是否启用" json:"enabled"`
}

func (Menu) TableName() string {
	return "menu"
}
