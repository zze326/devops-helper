package o_ci

import "github.com/zze326/devops-helper/app/modules/gormc"

type Env struct {
	gormc.Model
	Name          string `gorm:"comment:名称" json:"name"`
	Image         string `gorm:"comment:环境镜像" json:"image"`
	K8sSecretName string `gorm:"comment:K8s Secret Name" json:"k8s_secret_name"`
}

func (Env) TableName() string {
	return "ci_env"
}
