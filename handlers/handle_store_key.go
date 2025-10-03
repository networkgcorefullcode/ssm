package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

func HandleStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postStoreKey(mgr, w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

func postStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req models.StoreKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	label := req.KeyLabel
	id := req.Id
	key_value, err := base64.StdEncoding.DecodeString(req.KeyValue)
	if err != nil {
		http.Error(w, "bad base64", http.StatusBadRequest)
		return
	}

	if handle, err := mgr.StoreKey(label, key_value, []byte(id)); err != nil {
		http.Error(w, "failed to store key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := models.StoreKeyResponse{
		Handle:    handle,
		CipherKey: nil, // Puedes asignar el valor adecuado si tienes el cipher key
	}

	findHandle, err := mgr.FindKeyByLabel(constants.LABEL_ENCRYPTION_KEY)
	if err != nil || findHandle == 0 {
		json.NewEncoder(w).Encode(resp)
		return
	}

	cipher, err := mgr.EncryptKey(findHandle, nil, key_value, pkcs11.CKM_AES_CBC_PAD)
	if err != nil {
		resp.CipherKey = nil
		json.NewEncoder(w).Encode(resp)
		return
	}
	*resp.CipherKey = base64.StdEncoding.EncodeToString(cipher)
	json.NewEncoder(w).Encode(resp)

}
