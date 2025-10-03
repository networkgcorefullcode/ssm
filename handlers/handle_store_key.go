package handlers

import (
	"encoding/json"
	"net/http"

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
	var req models.StoreKey

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	label := req.KeyLabel
	id := req.ID
	key_value := req.KeyValue

	if handle, err := mgr.StoreKey(label, []byte(key_value), []byte(id)); err != nil {
		http.Error(w, "failed to store key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"handle":  handle,
		"message": "key stored successfully",
	})
}
