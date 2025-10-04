package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// HandleEncrypt maneja las peticiones de encriptación
// @Summary Encriptar datos
// @Description Encripta datos usando una clave AES almacenada en el HSM
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.EncryptRequest true "Datos a encriptar"
// @Success 200 {object} models.EncryptResponse "Datos encriptados exitosamente"
// @Failure 400 {object} models.ProblemDetails "Petición inválida"
// @Failure 404 {object} models.ProblemDetails "Clave no encontrada"
// @Failure 500 {object} models.ProblemDetails "Error interno del servidor"
// @Router /encrypt [post]
func HandleEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postEncrypt(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "El método HTTP no está permitido para este endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postEncrypt(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing encrypt request")

	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Decoding base64 plaintext for key: %s", req.KeyLabel)
	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode base64 plaintext: %v", err)
		sendProblemDetails(w, "Bad Request", "Los datos en base64 no son válidos", "INVALID_BASE64", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
	keyHandle, err := mgr.FindKeyByLabel(req.KeyLabel)
	if err != nil {
		logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(w, "Key Not Found", "La clave especificada no existe en el HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
		return
	}

	logger.AppLog.Info("Generating initialization vector (IV)")
	iv := make([]byte, 16)
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Error al generar el vector de inicialización", "IV_GENERATION_FAILED", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Encrypting data")
	ciphertext, err := mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC_PAD)
	safe.Zero(pt) // Limpiar datos sensibles de la memoria
	if err != nil {
		logger.AppLog.Errorf("Encryption failed: %v", err)
		sendProblemDetails(w, "Encryption Failed", "Error durante el proceso de encriptación", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Encryption completed successfully")

	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)
	ivStr := base64.StdEncoding.EncodeToString(iv)
	ok := true

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Crear respuesta usando el modelo estructurado
	resp := models.EncryptResponse{
		CipherB64:   &ciphertextStr,
		IvB64:       &ivStr,
		Ok:          &ok,
		TimeCreated: &timeCreated,
		TimeUpdated: &timeUpdated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}
