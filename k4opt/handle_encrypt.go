package k4opt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

func HandleEncryptK4(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Debugf("Received encryption request from %s", r.RemoteAddr)

	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	logger.AppLog.Debugf("Encryption request for key label: %s", req.KeyLabel)

	// Decodifica plaintext (base64 -> []byte)
	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode base64 plaintext: %v", err)
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

	logger.AppLog.Debugf("Decoded plaintext length: %d bytes", len(pt))

	handle := mgr.FindKeyByLabel(req.KeyLabel)

	if handle == 0 {
		logger.AppLog.Warnf("Key with label '%s' not found in HSM", req.KeyLabel)
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

	logger.AppLog.Debugf("Found key with handle: %d", handle)

	iv := make([]byte, 16)

	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		http.Error(w, "failed to generate IV", http.StatusInternalServerError)
		return
	}

	logger.AppLog.Debugf("Generated IV, starting encryption with AES key")

	// Use manager to encrypt
	ciphertext, err := mgr.EncryptWithAESKey(handle, iv, pt)

	// Scrub plain text
	safe.Zero(pt)
	logger.AppLog.Debug("Plaintext memory zeroed for security")

	if err != nil {
		logger.AppLog.Errorf("Encryption failed: %v", err)
		http.Error(w, "encrypt error", 500)
		return
	}

	logger.AppLog.Infof("Encryption successful for key '%s', ciphertext length: %d bytes", req.KeyLabel, len(ciphertext))

	now := time.Now().UTC().Format(time.RFC3339)
	resp := models.EncryptResponse{
		CipherB64:   base64.StdEncoding.EncodeToString(ciphertext),
		IVB64:       base64.StdEncoding.EncodeToString(iv),
		TimeCreated: now,
		TimeUpdated: now,
		Ok:          true,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		return
	}

	logger.AppLog.Debug("Encryption response sent successfully")
}
