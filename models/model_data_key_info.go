package models

type DataKeyInfo struct {
	// HSM key handle
	Handle int32 `json:"handle"`
	// Key identifier
	Id int32 `json:"id"`
}
