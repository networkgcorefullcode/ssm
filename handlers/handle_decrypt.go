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

func HandleDecrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		postDecrypt(mgr, w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}

}

func postDecrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req models.DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cipher, _ := base64.StdEncoding.DecodeString(req.CipherB64)
	iv, _ := base64.StdEncoding.DecodeString(req.IVB64)

	// Get handle using the label
	keyHandle, err := mgr.FindKeyByLabel(req.KeyLabel)
	if err != nil {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	// Get the plaintext using aes decrypt algoritm
	plaintext, err := mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC_PAD)
	if err != nil {
		http.Error(w, "decrypt error", 500)
		return
	}

	// At this point, the plaintext is in memory: use mlock before using it and zero afterwards
	// Send result (but ideally return only results, not the key)
	resp := map[string]string{
		"plain_b64": base64.StdEncoding.EncodeToString(plaintext),
	}
	// scrubbear
	safe.Zero(plaintext)
	_ = json.NewEncoder(w).Encode(resp)
}
