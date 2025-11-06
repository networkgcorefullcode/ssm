package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
)

// HandleHealthCheck handles health check requests
// @Summary Health check
// @Description Returns the health status of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health-check [get]
func HandleHealthCheck(c *gin.Context) {
	logger.AppLog.Info("Health check endpoint called")

	// Create response
	response := models.HealthCheckResponse{
		Status:  "OK",
		Message: "Service is healthy",
	}

	// Set response headers
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(c.Writer).Encode(response); err != nil {
		logger.AppLog.Errorf("Failed to encode health check response: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "Failed to encode response", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Health check response sent successfully")
}
