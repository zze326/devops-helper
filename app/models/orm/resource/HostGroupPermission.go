package o_resource

import (
	"gorm.io/gorm"
	"time"
)

type HostGroupPermission struct {
	ID         int            `gorm:"primarykey" json:"id,omitempty"`
	HostGroups []*HostGroup   `gorm:"many2many:host_group_permission_host_group;" json:"host_groups"`
	Type       int8           `gorm:"comment:类型 1:用户 2:角色;uniqueIndex:idx_type_ref_id" json:"type"`
	RefID      int            `gorm:"comment:关联 ID;uniqueIndex:idx_type_ref_id" json:"ref_id"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"uniqueIndex:idx_type_ref_id" json:"deleted_at,omitempty"`
}

func (HostGroupPermission) TableName() string {
	return "host_group_permission"
}
