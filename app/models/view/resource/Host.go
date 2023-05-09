package v_resource

type AddHostReq struct {
	GroupIDs                    []int  `json:"host_group_ids"`
	Name                        string `json:"name" valid:"required~名称不能为空"`
	Host                        string `json:"host" valid:"required~主机名不能为空"`
	Port                        int    `json:"port" valid:"int~端口必须为整型数字"`
	Username                    string `json:"username" valid:"required~用户名不能为空"`
	Password                    string `json:"password"`
	PrivateKey                  string `json:"private_key"`
	UseKey                      bool   `json:"use_key"`
	Desc                        string `json:"desc"`
	UpdatePasswordAndPrivateKey bool   `json:"update_password_and_private_key"`
}

type EditHostReq struct {
	ID int `json:"id" valid:"int~路由ID必须为整型数字,gt(0)~路由ID必须大于0"`
	AddHostReq
}
