package o_system

import "github.com/zze326/devops-helper/app/modules/gormc"

type User struct {
	gormc.Model
	Username string  `gorm:"comment:用户名" json:"username"`
	Password string  `gorm:"comment:密码" json:"password"`
	Phone    string  `gorm:"comment:手机号码" json:"phone"`
	Email    string  `gorm:"comment:邮箱" json:"email"`
	RealName string  `gorm:"comment:真实姓名" json:"real_name"`
	Roles    []*Role `gorm:"many2many:user_role;comment:角色" json:"roles"`
}

func (User) TableName() string {
	return "user"
}

func (u User) IsSuper() bool {
	if u.Username == "admin" {
		return true
	}
	for _, role := range u.Roles {
		if role.Code == "admin" {
			return true
		}
	}
	return false
}
