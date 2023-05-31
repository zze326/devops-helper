package o_resource

import (
	"github.com/zze326/devops-helper/app/consts"
	"github.com/zze326/devops-helper/app/modules/gormc"
)

// Secret 秘钥
type Secret struct {
	gormc.Model
	Name     string `gorm:"comment:名称" json:"name"`
	Type     int8   `gorm:"comment:类型" json:"type"`
	Content  string `gorm:"comment:内容;type:text" json:"content"`
	Username string `gorm:"comment:用户名" json:"username"`
	Password string `gorm:"comment:密码" json:"password"`
}

func (Secret) TableName() string {
	return "secret"
}

func (s *Secret) UseContent() bool {
	switch s.Type {
	case consts.SecretTypeKubernetesAuth:
		return true
	}
	return false
}
