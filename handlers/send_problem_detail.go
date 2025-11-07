package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/models"
)

// sendProblemDetails envía una respuesta de error usando ProblemDetails RFC 7807
// func sendProblemDetails(w http.ResponseWriter, title, detail, errorCode string, status int, instance string) {
// 	problem := models.ProblemDetails{
// 		Title:    title,
// 		Detail:   detail,
// 		Error:    errorCode,
// 		Status:   int32(status),
// 		Instance: instance,
// 	}

// 	w.Header().Set("Content-Type", "application/problem+json")
// 	w.WriteHeader(status)

// 	if err := json.NewEncoder(w).Encode(problem); err != nil {
// 		logger.AppLog.Errorf("Failed to encode problem details: %v", err)
// 	}
// }

// sendProblemDetails sends an RFC7807 problem+json error via Gin context
func sendProblemDetails(c *gin.Context, title, detail, errorCode string, status int, instance string) {
	problem := models.ProblemDetails{
		Title:    title,
		Detail:   detail,
		Error:    errorCode,
		Status:   int32(status),
		Instance: instance,
	}
	c.Error(errors.New(errorCode))
	c.JSON(status, problem)
}
