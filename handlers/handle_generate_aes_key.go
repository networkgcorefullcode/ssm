package handlers

import (
	"encoding/json"
	"net/http"

	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleGenerateAESKey maneja las peticiones de generación de claves AES
// @Summary Generar clave AES
// @Description Genera una nueva clave AES y la almacena en el HSM
// @Tags Key Management
// @Accept json
// @Produce json
// @Param request body models.GenAESKeyRequest true "Parámetros para generar la clave AES"
// @Success 201 {object} models.GenAESKeyResponse "Clave AES generada exitosamente"
// @Failure 400 {object} models.ProblemDetails "Petición inválida"
// @Failure 500 {object} models.ProblemDetails "Error interno del servidor"
// @Router /generate-aes-key [post]
func HandleGenerateAESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postGenerateAESKey(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "El método HTTP no está permitido para este endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postGenerateAESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing AES key generation request")

	var req models.GenAESKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}
	if req.Id <= 0 {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, r.URL.Path)
		return
	}
	if req.Bits != 128 && req.Bits != 256 {
		logger.AppLog.Errorf("Invalid key size: %d bits", req.Bits)
		sendProblemDetails(w, "Bad Request", "El tamaño de clave debe ser 128, 192 o 256 bits", "INVALID_KEY_SIZE", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Generating AES key - ID: %s, Bits: %d", req.Id, req.Bits)

	var label string
	if req.Bits == 128 {
		label = constants.LABEL_ENCRYPTION_KEY_AES128
	} else if req.Bits == 256 {
		label = constants.LABEL_ENCRYPTION_KEY_AES256
	}

	handle, err := mgr.GenerateAESKey(label, req.Id, int(req.Bits))
	if err != nil {
		logger.AppLog.Errorf("AES key generation failed: %v", err)
		sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave AES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("AES key generated successfully - Handle: %d", handle)

	resp := models.GenAESKeyResponse{
		Handle: int32(handle),
		Id:     req.Id,
		Bits:   req.Bits,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}
