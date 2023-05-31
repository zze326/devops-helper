package o_ci

import "github.com/zze326/devops-helper/app/modules/gormc"

type EnvRef struct {
	gormc.Model
	EnvID      int      `gorm:"comment:环境 ID" json:"env_id"`
	Env        *Env     `gorm:"foreignKey:EnvID" json:"env"`
	PipelineID int      `gorm:"comment:流水线 ID" json:"pipeline_id"`
	Sort       int      `gorm:"comment:排序" json:"sort"`
	Stages     []*Stage `gorm:"foreignKey:EnvRefID;preload:Task;preload:Tasks" json:"stages"`
}

func (EnvRef) TableName() string {
	return "ci_env_ref"
}
