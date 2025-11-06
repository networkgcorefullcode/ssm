package handlers

import (
	"encoding/hex"
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
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidJSON, ErrorCodeInvalidJSON, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Get MongoDB client from context or initialize
	client := database.Client

	// Query user from MongoDB
	filter := bson.M{"service_id": loginReq.ServiceId}
	userData, err := database.FindOneData(client, factory.SsmConfig.Configuration.Mongodb.DBName, database.CollSecret, filter)
	if err != nil {
		logger.AppLog.Errorf("User not found: %v", err)
		sendProblemDetails(c, ErrorTitleUnauthorized, "Invalid service ID or password", ErrorCodeUnauthorized, http.StatusUnauthorized, c.Request.URL.Path)
		return
	}

	// Decodify using bson
	user := database.UserSecret{}
	bsonBytes, _ := bson.Marshal(userData)
	if err := bson.Unmarshal(bsonBytes, &user); err != nil {
		logger.AppLog.Errorf("Failed to unmarshal user data: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "User data processing error", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	var iv []byte
	if iv, err = hex.DecodeString(user.PasswordSecret.IV); err != nil {
		logger.AppLog.Errorf("Failed to decode IV: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, ErrorDetailInvalidHexIV, ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Get PKCS11 session
	session := mgr.GetSession()
	defer mgr.LogoutSession(session)

	// Get key handle
	keyHandle, err := pkcs11mgr.FindKey(user.PasswordSecret.KeyLabel, user.PasswordSecret.Id, *session)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "Key retrieval error", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Decrypt stored password
	encryptedPassword, err := hex.DecodeString(user.PasswordSecret.EncryptedData)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode encrypted password: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "Password decryption error", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	decryptedPassword, err := pkcs11mgr.DecryptKey(keyHandle, iv, encryptedPassword, pkcs11.CKM_AES_CBC_PAD, *session)
	if err != nil {
		logger.AppLog.Errorf("Failed to decrypt password: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "Password decryption failed", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	// Compare passwords
	if hex.EncodeToString(decryptedPassword) != loginReq.Password {
		logger.AppLog.Warnf("Password mismatch for user: %s", loginReq.ServiceId)
		sendProblemDetails(c, ErrorTitleUnauthorized, "Invalid service ID or password", ErrorCodeUnauthorized, http.StatusUnauthorized, c.Request.URL.Path)
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ServiceID, session)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate JWT token: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, "Token generation failed", ErrorCodeInternalError, http.StatusInternalServerError, c.Request.URL.Path)
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
