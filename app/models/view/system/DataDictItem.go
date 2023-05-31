package v_system

type AddDataDictItemReq struct {
	Label      string `json:"label" validate:"required"`
	Value      int    `json:"value" validate:"required"`
	Sort       int    `json:"sort" validate:"required"`
	DataDictID int    `json:"data_dict_id" validate:"required"`
}

type EditDataDictItemReq struct {
	ID    int    `json:"id" validate:"required"`
	Label string `json:"label" validate:"required"`
	Value int    `json:"value" validate:"required"`
	Sort  int    `json:"sort" validate:"required"`
}
