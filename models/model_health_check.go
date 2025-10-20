package models

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
