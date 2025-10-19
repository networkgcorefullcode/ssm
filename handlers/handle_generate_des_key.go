package handlers

import (
	"encoding/json"
	"net/http"

	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleGenerateDESKey maneja las peticiones de generación de claves DES
// @Summary Generar clave DES
// @Description Genera una nueva clave DES y la almacena en el HSM
// @Tags Key Management
// @Accept json
// @Produce json
// @Param request body models.GenDESKeyRequest true "Parámetros para generar la clave DES"
// @Success 201 {object} models.GenDESKeyResponse "Clave DES generada exitosamente"
// @Failure 400 {object} models.ProblemDetails "Petición inválida"
// @Failure 500 {object} models.ProblemDetails "Error interno del servidor"
// @Router /generate-des-key [post]
func HandleGenerateDESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postGenerateDESKey(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "El método HTTP no está permitido para este endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postGenerateDESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing DES key generation request")

	var req models.GenDESKeyRequest
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

	logger.AppLog.Infof("Generating DES key - ID: %d", req.Id)
	handle, err := mgr.GenerateDESKey(constants.LABEL_K4_KEY_DES, req.Id)
	if err != nil {
		logger.AppLog.Errorf("DES key generation failed: %v", err)
		sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave DES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("DES key generated successfully - Handle: %d", handle)

	resp := models.GenDESKeyResponse{
		Handle: int32(handle),
		Id:     req.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}
