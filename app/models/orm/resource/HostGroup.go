package o_resource

import "github.com/zze326/devops-helper/app/modules/gormc"

type HostGroup struct {
	gormc.Model
	Name                 string                 `gorm:"comment:名称" json:"name"`
	ParentID             int                    `gorm:"comment:父级 ID" json:"parent_id"`
	Children             []*HostGroup           `gorm:"foreignKey:ParentID" json:"children"`
	Hosts                []*Host                `gorm:"many2many:host_group_host;" json:"hosts"`
	HostGroupPermissions []*HostGroupPermission `gorm:"many2many:host_group_permission_host_group;" json:"host_group_permissions"`
}

func (HostGroup) TableName() string {
	return "host_group"
}
