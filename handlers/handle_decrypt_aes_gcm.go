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
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidJSON, ErrorCodeInvalidJSON, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Debugf("Decryption request for key label: %s", req.KeyLabel)

	// Validate required fields
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, ErrorDetailKeyLabelRequired, ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Cipher == "" {
		logger.AppLog.Error("Ciphertext is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, "Ciphertext is required", ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Iv == "" {
		logger.AppLog.Error("IV is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, "IV is required for AES-GCM", ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Tag == "" {
		logger.AppLog.Error("Authentication tag is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, "Authentication tag is required for AES-GCM", ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Decode ciphertext, IV, and tag
	cipher, err := hex.DecodeString(req.Cipher)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode ciphertext hex: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexCiphertext, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	iv, err := hex.DecodeString(req.Iv)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode IV hex: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexIV, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	tag, err := hex.DecodeString(req.Tag)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode tag hex: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexTag, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Validate tag length (should be 16 bytes for 128-bit tag)
	if len(tag) != 16 {
		logger.AppLog.Errorf("Invalid tag length: %d bytes (expected 16)", len(tag))
		sendProblemDetails(c, ErrorTitleValidationError, "Authentication tag must be 16 bytes (128 bits)", ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Decode AAD if provided
	var aad []byte
	if req.Aad != "" {
		aad, err = hex.DecodeString(req.Aad)
		if err != nil {
			logger.AppLog.Errorf("Failed to decode hex AAD: %v", err)
			sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexAAD, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
			return
		}
		logger.AppLog.Infof("AAD provided: %d bytes", len(aad))
	}

	logger.AppLog.Debugf("Decoded ciphertext length: %d bytes, IV length: %d bytes, Tag length: %d bytes", len(cipher), len(iv), len(tag))

	// Find key by label
	keyHandle, err := pkcs11mgr.FindKey(req.KeyLabel, req.Id, *s)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key by label '%s': %v", req.KeyLabel, err)
		sendProblemDetails(c, ErrorTitleKeyNotFound, ErrorDetailKeyNotExist, ErrorCodeKeyNotFound, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Decrypting data with AES-GCM")

	// Combine ciphertext and tag for decryption (as expected by PKCS#11)
	ciphertextWithTag := append(cipher, tag...)

	// Decrypt with AES-GCM
	rawPlaintext, err := pkcs11mgr.DecryptKeyAesGCM(keyHandle, iv, ciphertextWithTag, aad, *s)
	if err != nil {
		logger.AppLog.Errorf("AES-GCM decryption failed (authentication may have failed): %v", err)
		sendProblemDetails(c, ErrorTitleDecryptionFailed, "AES-GCM decryption failed. The authentication tag may be invalid or the data may have been tampered with.", ErrorCodeDecryptionError, http.StatusUnauthorized, c.Request.URL.Path)
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
		Id:          req.Id,
	}

	// Clear sensitive data
	safe.Zero(rawPlaintext)

	c.JSON(http.StatusOK, resp)
}
