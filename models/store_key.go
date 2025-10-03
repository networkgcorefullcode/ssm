package models

type StoreKeyRequest struct {
	KeyLabel string `json:"key_label"`
	ID       string `json:"id"`
	KeyValue string `json:"key_value"`
}

type StoreKeyResponse struct {
	Handle    uint   `json:"handle"`
	CipherKey string `json:"cipher_key"`
}
