package v_system

type AddDataDictReq struct {
	Name     string `json:"name" validate:"required"`
	TypeCode string `json:"type_code" validate:"required"`
}

type EditDataDictReq struct {
	ID int `json:"id" validate:"required"`
	AddDataDictReq
}
