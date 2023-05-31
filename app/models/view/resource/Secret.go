package v_resource

type AddSecretReq struct {
	Name     string `json:"name" validate:"required"`
	Type     int8   `json:"type" validate:"required"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditSecretReq struct {
	AddSecretReq
	ID int `json:"id" validate:"required"`
}
