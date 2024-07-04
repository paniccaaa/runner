package models

type SharedCode struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	Output    string `json:"output"`
	ErrOutput string `json:"err_output"`
}
