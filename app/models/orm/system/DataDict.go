package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type DataDict struct {
	gormc.Model
	Name     string `gorm:"comment:名称" json:"name"`
	TypeCode string `gorm:"comment:类型代码;size:24" json:"type_code"`
}

func (DataDict) TableName() string {
	return "data_dict"
}
