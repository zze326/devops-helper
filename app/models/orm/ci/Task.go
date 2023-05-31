package o_ci

import "github.com/zze326/devops-helper/app/modules/gormc"

type Task struct {
	gormc.Model
	Type    int8   `gorm:"comment:任务类型" json:"type"`
	Url     string `gorm:"comment:地址" json:"url"`
	Branch  string `gorm:"comment:分支" json:"branch"`
	Content string `gorm:"comment:内容" json:"content"`
	StageID int    `gorm:"comment:阶段 ID" json:"stage_id"`
}

func (Task) TableName() string {
	return "ci_task"
}
