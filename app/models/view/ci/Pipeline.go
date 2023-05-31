package v_ci

type AddPipelineReq struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
}

type EditPipelineReq struct {
	ID int `json:"id" validate:"required"`
	AddPipelineReq
}
