package dto

type Task struct {
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Owner       string `json:"owner"`
	IsDone      bool   `json:"isDone"`
}
