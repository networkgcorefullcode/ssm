package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

func HandleEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postEncrypt(mgr, w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

func postEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		http.Error(w, "bad base64", http.StatusBadRequest)
		return
	}

	keyHandle, err := mgr.FindKeyByLabel(req.KeyLabel)
	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	iv := make([]byte, 16)
	if err := safe.RandRead(iv); err != nil {
		http.Error(w, "failed to generate IV", http.StatusInternalServerError)
		return
	}

	ciphertext, err := mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC_PAD)
	safe.Zero(pt)
	if err != nil {
		http.Error(w, "encrypt error", 500)
		return
	}

	resp := map[string]string{
		"cipher_b64": base64.StdEncoding.EncodeToString(ciphertext),
		"iv_b64":     base64.StdEncoding.EncodeToString(iv),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
