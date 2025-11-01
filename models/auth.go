package models

// LoginRequest represents the login request payload
type LoginRequest struct {
	ServiceID string `json:"service_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginResponse represents the JWT token response
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message,omitempty"`
}
