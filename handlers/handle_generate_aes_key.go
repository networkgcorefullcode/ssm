package handlers

import (
	"encoding/json"
	"net/http"

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

	// Validar parámetros requeridos
	if req.Label == "" {
		logger.AppLog.Error("Label is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'label' es requerido y no puede estar vacío", "MISSING_LABEL", http.StatusBadRequest, r.URL.Path)
		return
	}
	if req.Id == "" {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, r.URL.Path)
		return
	}
	if req.Bits != 128 && req.Bits != 192 && req.Bits != 256 {
		logger.AppLog.Errorf("Invalid key size: %d bits", req.Bits)
		sendProblemDetails(w, "Bad Request", "El tamaño de clave debe ser 128, 192 o 256 bits", "INVALID_KEY_SIZE", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Generating AES key - Label: %s, ID: %s, Bits: %d", req.Label, req.Id, req.Bits)
	handle, err := mgr.GenerateAESKey(req.Label, []byte(req.Id), req.Bits)
	if err != nil {
		logger.AppLog.Errorf("AES key generation failed: %v", err)
		sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave AES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("AES key generated successfully - Handle: %d", handle)

	resp := models.GenAESKeyResponse{
		Handle: handle,
		Label:  &req.Label,
		Id:     &req.Id,
		Bits:   &req.Bits,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}
