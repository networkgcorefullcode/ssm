package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
)

// sendProblemDetails envía una respuesta de error usando ProblemDetails RFC 7807
func sendProblemDetails(w http.ResponseWriter, title, detail, errorCode string, status int, instance string) {
	problem := models.ProblemDetails{
		Title:    title,
		Detail:   detail,
		Error:    errorCode,
		Status:   int32(status),
		Instance: instance,
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(problem); err != nil {
		logger.AppLog.Errorf("Failed to encode problem details: %v", err)
	}
}
