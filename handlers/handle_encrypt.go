package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// HandleEncrypt handles encryption requests
// @Summary Encrypt data
// @Description Encrypts data using an AES key stored in the HSM
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.EncryptRequest true "Data to encrypt"
// @Success 200 {object} models.EncryptResponse "Data encrypted successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 404 {object} models.ProblemDetails "Key not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /encrypt [post]
func HandleEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postEncrypt(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing encrypt request")

	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Decoding hex plaintext for key: %s", req.KeyLabel)
	pt, err := hex.DecodeString(req.Plain)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex plaintext: %v", err)
		sendProblemDetails(w, "Bad Request", "The hex data is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
	keyHandle, err := mgr.FindKey(req.KeyLabel, 0)
	if err != nil {
		logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(w, "Key Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
		return
	}

	logger.AppLog.Info("Generating initialization vector (IV)")
	iv := make([]byte, 16)
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Error generating initialization vector", "IV_GENERATION_FAILED", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Encrypting data")

	var encryptionAlgorithm uint
	switch req.EncryptionAlgorithm {
	case constants.ALGORITM_AES_128, constants.ALGORITM_AES_256:
		encryptionAlgorithm = pkcs11.CKM_AES_CBC_PAD
	case constants.ALGORITM_DES:
		encryptionAlgorithm = pkcs11.CKM_DES_CBC_PAD
	case constants.ALGORITM_DES3:
		encryptionAlgorithm = pkcs11.CKM_DES3_CBC_PAD
	}

	ciphertext, err := mgr.EncryptKey(keyHandle, iv, pt, encryptionAlgorithm)
	safe.Zero(pt) // Clear sensitive data from memory
	if err != nil {
		logger.AppLog.Errorf("Encryption failed: %v", err)
		sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Encryption completed successfully")

	ciphertextStr := hex.EncodeToString(ciphertext)
	ivStr := hex.EncodeToString(iv)
	ok := true

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Create response using the structured model
	resp := models.EncryptResponse{
		Cipher:      ciphertextStr,
		Iv:          ivStr,
		Ok:          ok,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}
