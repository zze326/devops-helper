package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type FrontendRoute struct {
	gormc.Model
	Name          string `gorm:"comment:名称" json:"name"`
	Title         string `gorm:"comment:标题" json:"title"`
	RoutePath     string `gorm:"comment:路由路径" json:"route_path"`
	ComponentPath string `gorm:"comment:组件路径" json:"component_path"`
}

func (FrontendRoute) TableName() string {
	return "frontend_route"
}
