package k4opt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

func HandleEncryptK4(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Decodifica plaintext (base64 -> []byte)
	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		problem := models.ProblemDetails{
			Title:  "Bad Base64",
			Status: http.StatusBadRequest,
			Err:    "bad_base64",
			Detail: "Failed to decode base64 plaintext",
		}
		json.NewEncoder(w).Encode(problem)
		return
	}

	handle := mgr.FindKeyByLabel(req.KeyLabel)

	if handle == 0 {
		w.WriteHeader(http.StatusNotFound)
		problem := models.ProblemDetails{
			Title:  "Key Not Found",
			Status: http.StatusNotFound,
			Err:    "key_not_found",
			Detail: fmt.Sprintf("Key with label '%s' not found", req.KeyLabel),
		}
		json.NewEncoder(w).Encode(problem)
		return
	}

	iv := make([]byte, 16)

	if err := safe.RandRead(iv); err != nil {
		http.Error(w, "failed to generate IV", http.StatusInternalServerError)
		return
	}

	// Use manager to encrypt
	ciphertext, err := mgr.EncryptWithAESKey(1, iv, pt)

	// Scrub plain text
	safe.Zero(pt)

	if err != nil {
		http.Error(w, "encrypt error", 500)
		return
	}

	resp := models.EncryptResponse{
		CipherB64:   base64.StdEncoding.EncodeToString(ciphertext),
		IVB64:       base64.StdEncoding.EncodeToString(iv),
		TimeCreated: req.KeyLabel,
		TimeUpdated: req.KeyLabel,
		Ok:          true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
