package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// HandleDecryptAESGCM handles AES-GCM decryption requests
// @Summary Decrypt data with AES-GCM
// @Description Decrypts data using an AES key stored in the HSM with GCM mode (authenticated decryption)
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.DecryptAESGCMRequest true "Data to decrypt with AES-GCM"
// @Success 200 {object} models.DecryptAESGCMResponse "Data decrypted successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 401 {object} models.ProblemDetails "Authentication failed (invalid tag)"
// @Failure 404 {object} models.ProblemDetails "Key not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /crypto/decrypt-aes-gcm [post]
func HandleDecryptAESGCM(c *gin.Context) {
	logger.AppLog.Info("Processing AES-GCM decrypt request")

	// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	var req models.DecryptAESGCMRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, "Invalid JSON", "Failed to parse request body: "+err.Error(), "INVALID_JSON", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Debugf("Decryption request for key label: %s", req.KeyLabel)

	// Validate required fields
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(c, "Validation Error", "Key label is required", "VALIDATION_FAILED", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Cipher == "" {
		logger.AppLog.Error("Ciphertext is required but was empty")
		sendProblemDetails(c, "Validation Error", "Ciphertext is required", "VALIDATION_FAILED", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Iv == "" {
		logger.AppLog.Error("IV is required but was empty")
		sendProblemDetails(c, "Validation Error", "IV is required for AES-GCM", "VALIDATION_FAILED", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Tag == "" {
		logger.AppLog.Error("Authentication tag is required but was empty")
		sendProblemDetails(c, "Validation Error", "Authentication tag is required for AES-GCM", "VALIDATION_FAILED", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Decode ciphertext, IV, and tag
	cipher, err := hex.DecodeString(req.Cipher)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode ciphertext hex: %v", err)
		sendProblemDetails(c, "Invalid hex", "Failed to decode ciphertext: "+err.Error(), "INVALID_HEX", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	iv, err := hex.DecodeString(req.Iv)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode IV hex: %v", err)
		sendProblemDetails(c, "Invalid hex", "Failed to decode IV: "+err.Error(), "INVALID_HEX", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	tag, err := hex.DecodeString(req.Tag)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode tag hex: %v", err)
		sendProblemDetails(c, "Invalid hex", "Failed to decode authentication tag: "+err.Error(), "INVALID_HEX", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Validate tag length (should be 16 bytes for 128-bit tag)
	if len(tag) != 16 {
		logger.AppLog.Errorf("Invalid tag length: %d bytes (expected 16)", len(tag))
		sendProblemDetails(c, "Validation Error", "Authentication tag must be 16 bytes (128 bits)", "INVALID_TAG_LENGTH", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Decode AAD if provided
	var aad []byte
	if req.Aad != "" {
		aad, err = hex.DecodeString(req.Aad)
		if err != nil {
			logger.AppLog.Errorf("Failed to decode hex AAD: %v", err)
			sendProblemDetails(c, "Bad Request", "The AAD hex data is not valid", "INVALID_HEX", http.StatusBadRequest, c.Request.URL.Path)
			return
		}
		logger.AppLog.Infof("AAD provided: %d bytes", len(aad))
	}

	logger.AppLog.Debugf("Decoded ciphertext length: %d bytes, IV length: %d bytes, Tag length: %d bytes", len(cipher), len(iv), len(tag))

	// Find key by label
	keyHandle, err := pkcs11mgr.FindKeyLabelReturnRandom(req.KeyLabel, *s)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key by label '%s': %v", req.KeyLabel, err)
		sendProblemDetails(c, "Key Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, c.Request.URL.Path)
		return
	}

	// Get key attributes
	attr, err := pkcs11mgr.GetObjectAttributes(keyHandle, *s)
	if err != nil {
		logger.AppLog.Errorf("Attributes not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(c, "Attributes Not Found", "Failed to retrieve key attributes", "ATTRIBUTES_NOT_FOUND", http.StatusNotFound, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Decrypting data with AES-GCM")

	// Combine ciphertext and tag for decryption (as expected by PKCS#11)
	ciphertextWithTag := append(cipher, tag...)

	// Decrypt with AES-GCM
	rawPlaintext, err := pkcs11mgr.DecryptKeyAesGCM(keyHandle, iv, ciphertextWithTag, aad, *s)
	if err != nil {
		logger.AppLog.Errorf("AES-GCM decryption failed (authentication may have failed): %v", err)
		sendProblemDetails(c, "Decryption Failed", "AES-GCM decryption failed. The authentication tag may be invalid or the data may have been tampered with.", "DECRYPTION_ERROR", http.StatusUnauthorized, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("AES-GCM decryption completed successfully")

	// Convert to hex
	plainHex := hex.EncodeToString(rawPlaintext)

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Create response using the structured model
	resp := models.DecryptAESGCMResponse{
		Plain:       plainHex,
		Ok:          true,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
		Id:          attr.Id,
	}

	// Clear sensitive data
	safe.Zero(rawPlaintext)

	c.JSON(http.StatusOK, resp)
}
