package models

type GenAESKeyRequest struct {
	Label string `json:"label"`
	ID    string `json:"id"`
	Bits  int    `json:"bits"`
}

type GenAESKeyResponse struct {
	Handle uint   `json:"handle"`
	Label  string `json:"label"`
	ID     string `json:"id"`
	Bits   int    `json:"bits"`
}
