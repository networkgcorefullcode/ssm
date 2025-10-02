package k4opt

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

type decryptRequest struct {
	KeyLabel  string `json:"key_label"`
	CipherB64 string `json:"cipher_b64"`
	IVB64     string `json:"iv_b64"`
}

func HandleDecryptK4(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req decryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cipher, _ := base64.StdEncoding.DecodeString(req.CipherB64)
	iv, _ := base64.StdEncoding.DecodeString(req.IVB64)

	// Obtener handle por label (implementar)
	plaintext, err := mgr.DecryptWithAESKey( /*handle*/ 1, iv, cipher)
	if err != nil {
		http.Error(w, "decrypt error", 500)
		return
	}
	// A estas alturas, el plaintext está en memoria: usar mlock antes de usarlo y zero luego
	// Enviar result (pero idealmente devolver solo resultados, no la clave)
	resp := map[string]string{
		"plain_b64": base64.StdEncoding.EncodeToString(plaintext),
	}
	// scrubbear
	safe.Zero(plaintext)
	_ = json.NewEncoder(w).Encode(resp)
}
