package v_system

type AddRoleReq struct {
	Name          string `json:"name" valid:"required~角色名称不能为空"`
	Code          string `json:"code" valid:"required~角色代码不能为空"`
	PermissionIDs []int  `json:"permission_ids"`
}

type EditRoleReq struct {
	ID int `json:"id" valid:"int~角色ID必须为整型数字,gt(0)~角色ID必须大于0"`
	AddRoleReq
}
