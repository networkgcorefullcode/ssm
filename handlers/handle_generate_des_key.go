package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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
func HandleGenerateDESKey(c *gin.Context) {
	logger.AppLog.Info("Processing DES key generation request")

	// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	var req models.GenDESKeyRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, c.Request.URL.Path)
		return
	}
	if req.Id <= 0 {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(c, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("Generating DES key - ID: %d", req.Id)
	handle, err := pkcs11mgr.GenerateDESKey(constants.LABEL_ENCRYPTION_KEY_DES, req.Id, *s)
	if err != nil {
		logger.AppLog.Errorf("DES key generation failed: %v", err)
		sendProblemDetails(c, "Key Generation Failed", "Error al generar la clave DES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("DES key generated successfully - Handle: %d", handle)

	resp := models.GenDESKeyResponse{
		Handle: int32(handle),
		Id:     req.Id,
	}

	c.JSON(http.StatusCreated, resp)
}
