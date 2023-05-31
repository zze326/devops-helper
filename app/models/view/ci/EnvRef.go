package v_ci

type SaveEnvRefsReq struct {
	PipelineID int       `json:"pipeline_id" validate:"required"`
	EnvRefs    []*EnvRef `json:"env_refs" validate:"required"`
}

type EnvRef struct {
	ID     int      `json:"id"`
	EnvID  int      `json:"env_id" validate:"required"`
	Sort   int      `json:"sort" validate:"required"`
	Stages []*Stage `json:"stages"`
}

type Stage struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" validate:"required"`
	Parallel bool    `json:"parallel"`
	Sort     int     `json:"sort" validate:"required"`
	Task     *Task   `json:"task"`
	Tasks    []*Task `json:"tasks"`
}

type Task struct {
	ID      int    `json:"id"`
	Type    int8   `json:"type"`
	Content string `json:"content"`
	Url     string `json:"url"`
	Branch  string `json:"branch"`
}
