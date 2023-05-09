package v_system

type AddBackendRouteReq struct {
	Name string `json:"name" valid:"required~路由名称不能为空"`
	Path string `json:"path" valid:"required~路由路径不能为空"`
}

type EditBackendRouteReq struct {
	ID int `json:"id,optional" valid:"int~路由ID必须为整型数字,gt(0)~路由ID必须大于0"`
	AddBackendRouteReq
}
