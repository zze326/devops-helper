package v_system

type AddMenuReq struct {
	Name     string `json:"name" valid:"required~菜单名称不能为空"`
	Icon     string `json:"icon" valid:"stringlength(3|18)~图标名称长度为3-12位"`
	RouteID  int    `json:"route_id"`
	ParentID int    `json:"parent_id"`
	Sort     int    `json:"sort"`
	Enabled  bool   `json:"enabled"`
}

type UpdateMenuReq struct {
	ID int `json:"id" valid:"int~菜单ID必须为整型数字,gt(0)~菜单ID必须大于0"`
	AddMenuReq
}
