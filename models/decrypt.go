package models

type DecryptRequest struct {
	KeyLabel  string `json:"key_label"`
	CipherB64 string `json:"cipher_b64"`
	IVB64     string `json:"iv_b64"`
	ID        int    `json:"id,omitempty"`
}

type DecryptResponse struct {
	PlainB64 string `json:"plain_b64"`
}
