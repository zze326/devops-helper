package o_ci

import "github.com/zze326/devops-helper/app/modules/gormc"

type Pipeline struct {
	gormc.Model
	Name    string    `gorm:"comment:名称" json:"name"`
	Desc    string    `gorm:"comment:描述" json:"desc"`
	EnvRefs []*EnvRef `gorm:"foreignKey:PipelineID" json:"ci_env_refs"`
}

func (Pipeline) TableName() string {
	return "ci_pipeline"
}
