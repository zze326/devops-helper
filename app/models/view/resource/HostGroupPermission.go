package v_resource

type AddHostGroupPermissionReq struct {
	HostGroupIDs []int `json:"host_group_ids" valid:"required~服务器组ID不能为空"`
	Type         int8  `json:"type" valid:"required~类型不能为空,int~类型必须为整型数字"`
	RefID        int   `json:"ref_id" valid:"required~关联ID不能为空,int~关联ID必须为整型数字"`
}

type EditHostGroupPermissionReq struct {
	ID int `json:"id" valid:"required~服务器组ID不能为空"`
	AddHostGroupPermissionReq
}

type ListPageHostGroupPermissionItem struct {
	ID                int    `json:"id"`
	Type              int8   `json:"type"`
	RefID             int    `json:"ref_id"`
	RefName           string `json:"ref_name"`
	HostGroupNamesStr string `json:"host_group_names_str"`
}
