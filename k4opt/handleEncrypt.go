package k4opt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/networkgcorefullcode/ssm/safe"
)

type encryptRequest struct {
	KeyLabel string `json:"key_label"`
	PlainB64 string `json:"plain_b64"`
}

func HandleEncryptK4(w http.ResponseWriter, r *http.Request) {
	var req encryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Decodifica plaintext (base64 -> []byte)
	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		http.Error(w, "bad base64", http.StatusBadRequest)
		return
	}

	// Aquí deberías buscar el handle por etiqueta; ejemplo simplificado:
	// handle := s.mgr.FindKeyByLabel(req.KeyLabel)  // implementar
	// for demo, asumimos que ya tienes handle:
	// iv := random 16 bytes
	// iv := make([]byte, 16)
	// TODO: usar crypto/rand.Read(iv)
	// Usa manager para encrypt
	// ciphertext, err := s.mgr.EncryptWithAESKey( /*handle*/ 1, iv, pt)
	// Scrub plain text
	safe.Zero(pt)

	if err != nil {
		http.Error(w, "encrypt error", 500)
		return
	}

	// Almacena ciphertext + iv + metadata en Mongo
	// Ejemplo simplificado:

	fmt.Fprintf(w, `{"ok": true}`)
}
