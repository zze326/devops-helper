package v_system

import orm_system "github.com/zze326/devops-helper/app/models/orm/system"

type AddPermissionReq struct {
	Name             string `json:"name" valid:"required~权限名称不能为空"`
	Code             string `json:"code" valid:"required~权限代码不能为空"`
	ParentID         int    `json:"parent_id"`
	ChildrenIDs      []int  `json:"children_ids"`
	MenuIDs          []int  `json:"menu_ids"`
	FrontendRouteIDs []int  `json:"frontend_route_ids"`
	BackendRouteIDs  []int  `json:"backend_route_ids"`
}

type EditPermissionReq struct {
	ID int `json:"id" valid:"int~权限ID必须为整型数字,gt(0)~权限ID必须大于0"`
	AddPermissionReq
}

type GetPermissionResp struct {
	ID               int                      `json:"id"`
	Code             string                   `json:"code"`
	Name             string                   `json:"name"`
	ParentID         int                      `json:"parent_id"`
	Children         []*orm_system.Permission `json:"children"`
	ChildrenIDs      []int                    `json:"children_ids"`
	MenuIDs          []int                    `json:"menu_ids"`
	FrontendRouteIDs []int                    `json:"frontend_route_ids"`
	BackendRouteIDs  []int                    `json:"backend_route_ids"`
}
