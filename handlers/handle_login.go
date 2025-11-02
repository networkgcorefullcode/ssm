package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/database"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleLogin(c *gin.Context) {
	var loginReq models.LoginRequest

	// Parse JSON request
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		logger.AppLog.Errorf("Invalid JSON payload: %v", err)
		sendProblemDetails(c, "Bad Request", "Invalid JSON payload", "invalid_request", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Get MongoDB client from context or initialize
	client := database.Client

	// Query user from MongoDB
	filter := bson.M{"service_id": loginReq.ServiceId}
	userData, err := database.FindOneData(client, factory.SsmConfig.Configuration.Mongodb.DBName, database.CollSecret, filter)
	if err != nil {
		logger.AppLog.Errorf("User not found: %v", err)
		sendProblemDetails(c, "Unauthorized", "Invalid service ID or password", "invalid_credentials", http.StatusUnauthorized, c.Request.URL.Path)
		return
	}

	// Parse user data
	user := database.UserSecret{}

	// Convert userData to JSON and then unmarshal to user struct
	jsonData, err := json.Marshal(userData)
	if err != nil {
		logger.AppLog.Errorf("Failed to marshal user data: %v", err)
		sendProblemDetails(c, "Internal Server Error", "User data processing error", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	if err := json.Unmarshal(jsonData, &user); err != nil {
		logger.AppLog.Errorf("Failed to unmarshal user data: %v", err)
		sendProblemDetails(c, "Internal Server Error", "User data processing error", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	var iv []byte
	if iv, err = hex.DecodeString(user.PasswordSecret.IV); err != nil {
		logger.AppLog.Errorf("Failed to decode IV: %v", err)
		sendProblemDetails(c, "Internal Server Error", "Invalid IV format", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Get PKCS11 session
	session := mgr.GetSession()
	defer mgr.LogoutSession(session)

	// Get key handle
	keyHandle, err := pkcs11mgr.FindKey(user.PasswordSecret.KeyLabel, user.PasswordSecret.Id, *session)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key: %v", err)
		sendProblemDetails(c, "Internal Server Error", "Key retrieval error", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Decrypt stored password
	encryptedPassword, err := hex.DecodeString(user.PasswordSecret.EncryptedData)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode encrypted password: %v", err)
		sendProblemDetails(c, "Internal Server Error", "Password decryption error", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	decryptedPassword, err := pkcs11mgr.DecryptKey(keyHandle, iv, encryptedPassword, pkcs11.CKM_AES_CBC_PAD, *session)
	if err != nil {
		logger.AppLog.Errorf("Failed to decrypt password: %v", err)
		sendProblemDetails(c, "Internal Server Error", "Password decryption failed", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Compare passwords
	if string(decryptedPassword) != loginReq.Password {
		logger.AppLog.Warnf("Password mismatch for user: %s", loginReq.ServiceId)
		sendProblemDetails(c, "Unauthorized", "Invalid service ID or password", "invalid_credentials", http.StatusUnauthorized, c.Request.URL.Path)
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ServiceID, session)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate JWT token: %v", err)
		sendProblemDetails(c, "Internal Server Error", "Token generation failed", "internal_error", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("User %s logged in successfully", user.ServiceID)
	c.JSON(http.StatusOK, models.LoginResponse{
		Token:   token,
		Message: "Login successful",
	})
}

func generateJWTToken(username string, s *pkcs11mgr.Session) (string, error) {
	token, err := pkcs11mgr.CreateStandardJWT(s, "ssm-service", username, "ssm-clients", 24)
	if err != nil {
		return "", err
	}

	return token, nil
}
