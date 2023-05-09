package v_system

type AppLogin struct {
	Username string `json:"username" valid:"required~用户名不能为空,stringlength(3|12)~用户名长度为3-12位"`
	Password string `json:"password" valid:"required~密码不能为空,stringlength(6|24)~密码长度为6-24位"`
}
