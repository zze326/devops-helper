package v_system

type AddFrontendRouteReq struct {
	Name          string `json:"name" valid:"required~路由名称不能为空"`
	Title         string `json:"title" valid:"required~密码不能为空"`
	ComponentPath string `json:"component_path" valid:"required~组件路径不能为空"`
	RoutePath     string `json:"route_path" valid:"required~路由路径不能为空"`
}

type EditFrontendRouteReq struct {
	ID int `json:"id,optional" valid:"int~路由ID必须为整型数字,gt(0)~路由ID必须大于0"`
	AddFrontendRouteReq
}
