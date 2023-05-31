package o_ci

import "github.com/zze326/devops-helper/app/modules/gormc"

type Stage struct {
	gormc.Model
	Name     string  `gorm:"comment:名称" json:"name"`
	Parallel bool    `gorm:"comment:并行" json:"parallel"`
	Sort     int     `gorm:"comment:排序" json:"sort"`
	EnvRefID int     `gorm:"comment:环境引用 ID" json:"env_ref_id"`
	TaskID   int     `gorm:"comment:串行任务 ID" json:"task_id"`
	Task     *Task   `gorm:"foreignKey:TaskID" json:"task"`
	Tasks    []*Task `gorm:"foreignKey:StageID" json:"tasks"`
}

func (Stage) TableName() string {
	return "ci_stage"
}

//name: 拉取代码
//parallel: false
//task:
//plugin: git-checkout
//configuration:
//url: http://192.168.1.195:8990/scm/ltcs/loverent-activity-application.git
//branch: master
