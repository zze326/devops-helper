package v_resource

type AddHostGroupReq struct {
	ParentID int    `json:"parent_id" valid:"int~父级ID必须为整型数字"`
	Name     string `json:"name" valid:"required~名称不能为空"`
}

type RenameHostGroupReq struct {
	ID   int    `json:"id" valid:"required~ID不能为空,int~ID必须为整型数字"`
	Name string `json:"name" valid:"required~名称不能为空"`
}
