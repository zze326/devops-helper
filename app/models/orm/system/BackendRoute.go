package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type BackendRoute struct {
	gormc.Model
	Name string `gorm:"comment:名称" json:"name"`
	Path string `gorm:"comment:路径" json:"path"`
}

func (BackendRoute) TableName() string {
	return "backend_route"
}
