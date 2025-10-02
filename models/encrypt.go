package models

type EncryptRequest struct {
	KeyLabel string `json:"key_label"`
	PlainB64 string `json:"plain_b64"`
}

type EncryptResponse struct {
	CipherB64   string `json:"cipher_b64"`
	IVB64       string `json:"iv_b64"`
	Ok          bool   `json:"ok"`
	TimeCreated string `json:"time_created,omitempty"`
	TimeUpdated string `json:"time_updated,omitempty"`
}
