package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type Role struct {
	gormc.Model
	Name        string        `gorm:"type:varchar(48);comment:角色名称;" json:"name"`
	Code        string        `gorm:"type:varchar(48);comment:角色代码;" json:"code"`
	Permissions []*Permission `gorm:"many2many:role_permission;" json:"permissions"`
}

func (Role) TableName() string {
	return "role"
}
