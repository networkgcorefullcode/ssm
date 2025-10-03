package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

func HandleGenerateAESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postGenerateAESKey(mgr, w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

func postGenerateAESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	var req models.GenAESKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Label == "" {
		http.Error(w, "label is required", http.StatusBadRequest)
		return
	}
	if req.Id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if req.Bits != 128 && req.Bits != 192 && req.Bits != 256 {
		http.Error(w, "bits must be 128, 192, or 256", http.StatusBadRequest)
		return
	}

	handle, err := mgr.GenerateAESKey(req.Label, []byte(req.Id), req.Bits)
	if err != nil {
		http.Error(w, "key generation failed", http.StatusInternalServerError)
		return
	}

	resp := models.GenAESKeyResponse{
		Handle: uint(handle),
		Label:  &req.Label,
		Id:     &req.Id,
		Bits:   &req.Bits,
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
