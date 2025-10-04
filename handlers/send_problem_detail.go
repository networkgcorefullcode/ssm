package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
)

// sendProblemDetails envía una respuesta de error usando ProblemDetails RFC 7807
func sendProblemDetails(w http.ResponseWriter, title, detail, errorCode string, status int, instance string) {
	problem := models.NewProblemDetails()
	problem.SetTitle(title)
	problem.SetDetail(detail)
	problem.SetError(errorCode)
	problem.SetStatus(int32(status))
	problem.SetInstance(instance)

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(problem); err != nil {
		logger.AppLog.Errorf("Failed to encode problem details: %v", err)
	}
}
