package models

type ProblemDetails struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Err    string `json:"error,omitempty"`
	Status int    `json:"status,omitempty"`
}
