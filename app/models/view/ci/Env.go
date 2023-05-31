package v_ci

type AddEnvReq struct {
	Name          string `json:"name" validate:"required"`
	Image         string `json:"image" validate:"required"`
	K8sSecretName string `json:"k8s_secret_name"`
}

type EditEnvReq struct {
	ID int `json:"id" validate:"required"`
	AddEnvReq
}
