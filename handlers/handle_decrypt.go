package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// @title        Decrypt Data API
// @version 	 1.0.0
// @description  Decrypt a Key using simetrics algoritms as AES 128 and AES 256, DES, 3DES
// @Accept       json
// @Produce      json
// @Param        request  body      models.DecryptRequest  true  "Datos para descifrar"
// @Success      200      {object}  models.DecryptResponse "Datos descifrados exitosamente"
// @Failure      400      {object}  models.ProblemDetails  "Error de validación o JSON inválido"
// @Failure      404      {object}  models.ProblemDetails  "Clave no encontrada"
// @Failure      405      {object}  models.ProblemDetails  "Método HTTP no permitido"
// @Failure      500      {object}  models.ProblemDetails  "Error interno del servidor"
// @Router       /decrypt [post]
func HandleDecrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Debugf("Received decrypt request from %s", r.RemoteAddr)

	switch r.Method {
	case http.MethodPost:
		postDecrypt(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "Only POST method is allowed", "method_not_allowed", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postDecrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Debugf("Processing decrypt request for %s", r.URL.Path)

	var req models.DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Invalid JSON", "Failed to parse request body: "+err.Error(), "bad_json", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Debugf("Decryption request for key label: %s", req.KeyLabel)

	// Validar campos requeridos
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(w, "Validation Error", "Key label is required", "validation_failed", http.StatusBadRequest, r.URL.Path)
		return
	}

	if req.CipherB64 == "" {
		logger.AppLog.Error("Ciphertext is required but was empty")
		sendProblemDetails(w, "Validation Error", "Ciphertext is required", "validation_failed", http.StatusBadRequest, r.URL.Path)
		return
	}

	// Decode ciphertext y IV
	cipher, err := base64.StdEncoding.DecodeString(req.CipherB64)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode ciphertext Base64: %v", err)
		sendProblemDetails(w, "Invalid Base64", "Failed to decode ciphertext: "+err.Error(), "bad_base64", http.StatusBadRequest, r.URL.Path)
		return
	}

	iv, err := base64.StdEncoding.DecodeString(req.IvB64)

	if err != nil {
		logger.AppLog.Errorf("Failed to decode IV Base64: %v", err)
		iv = nil
	}
	if req.IvB64 == "" {
		logger.AppLog.Info("iv is empty")
		iv = nil
	}

	logger.AppLog.Debugf("Decoded ciphertext length: %d bytes, IV length: %d bytes", len(cipher), len(iv))

	// Find key by label
	keyHandle, err := mgr.FindKeyByLabel(req.KeyLabel)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key by label '%s': %v", req.KeyLabel, err)
		sendProblemDetails(w, "Key Search Failed", "Failed to search for key: "+err.Error(), "key_search_failed", http.StatusInternalServerError, r.URL.Path)
		return
	}

	var plaintext []byte
	switch req.EncryptionAlgoritme {
	case constants.ALGORITM_AES_128:
		// Get the plaintext using aes decrypt algoritm
		plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC_PAD)
	case constants.ALGORITM_AES_256:
		// Get the plaintext using aes decrypt algoritm
		plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC_PAD)
	case constants.ALGORITM_DES:
		// Get the plaintext using aes decrypt algoritm
		plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_CBC_PAD)
	case constants.ALGORITM_DES3:
		// Get the plaintext using aes decrypt algoritm
		plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_CBC_PAD)
	}

	if err != nil {
		logger.AppLog.Errorf("Decryption failed: %v", err)
		sendProblemDetails(w, "Decryption Failed", "Failed to decrypt data: "+err.Error(), "decryption_failed", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Decryption successful for key '%s', plaintext length: %d bytes", req.KeyLabel, len(plaintext))

	// Preparar respuesta
	w.Header().Set("Content-Type", "application/json")
	resp := models.DecryptResponse{
		PlainB64: nil,
	}

	resp.SetPlainB64(base64.StdEncoding.EncodeToString(plaintext))

	// Limpiar memoria del plaintext por seguridad
	safe.Zero(plaintext)
	logger.AppLog.Debug("Plaintext memory zeroed for security")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		sendProblemDetails(w, "Response Encoding Failed", "Failed to encode response: "+err.Error(), "encoding_failed", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Debug("Decryption response sent successfully")
}
