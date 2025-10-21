package handlers

import (
	"encoding/json"
	"net/http"

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
func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Health check endpoint called")

	// Create response
	response := models.HealthCheckResponse{
		Status:  "OK",
		Message: "Service is healthy",
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AppLog.Errorf("Failed to encode health check response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.AppLog.Info("Health check response sent successfully")
}
