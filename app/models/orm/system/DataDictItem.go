package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type DataDictItem struct {
	gormc.Model
	Label      string `gorm:"comment:标签" json:"label"`
	Value      int    `gorm:"comment:值" json:"value"`
	Sort       int    `gorm:"comment:排序" json:"sort"`
	DataDictID int    `gorm:"comment:数据字典 ID" json:"data_dict_id"`
}

func (DataDictItem) TableName() string {
	return "data_dict_item"
}
