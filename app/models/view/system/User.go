package v_system

type AddUserReq struct {
	Username string `json:"username" valid:"required~用户名不能为空,stringlength(3|12)~用户名长度为3-12位"`
	Phone    string `json:"phone" valid:"matches(^\\d{11}$)~手机号码格式错误"`
	Email    string `json:"email" valid:"email~邮箱格式错误"`
	RealName string `json:"real_name" valid:"required~真实姓名不能为空,runelength(2|4)~真实姓名长度为2-4位"`
	RoleIDs  []int  `json:"role_ids"`
}

type EditUserReq struct {
	ID int `json:"id" valid:"int~路由ID必须为整型数字,gt(0)~路由ID必须大于0"`
	AddUserReq
}

type ResetUserPasswordReq struct {
	ID          int    `json:"id" valid:"int~路由ID必须为整型数字,gt(0)~路由ID必须大于0"`
	NewPassword string `json:"new_password" valid:"required~新密码不能为空,stringlength(6|24)~新密码长度为6-24位"`
}
