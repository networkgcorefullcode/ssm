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

// HandleEncryptAESGCM handles AES-GCM encryption requests
// @Summary Encrypt data with AES-GCM
// @Description Encrypts data using an AES key stored in the HSM with GCM mode (authenticated encryption)
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.EncryptAESGCMRequest true "Data to encrypt with AES-GCM"
// @Success 200 {object} models.EncryptAESGCMResponse "Data encrypted successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 404 {object} models.ProblemDetails "Key not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /crypto/encrypt-aes-gcm [post]
func HandleEncryptAESGCM(c *gin.Context) {
	// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	logger.AppLog.Info("Processing AES-GCM encrypt request")

	var req models.EncryptAESGCMRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidJSON, ErrorCodeInvalidJSON, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Validate required fields
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, ErrorDetailKeyLabelRequired, ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Plain == "" {
		logger.AppLog.Error("Plaintext is required but was empty")
		sendProblemDetails(c, ErrorTitleValidationError, ErrorDetailPlaintextRequired, ErrorCodeValidationFailed, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("Decoding hex plaintext for key: %s", req.KeyLabel)
	pt, err := hex.DecodeString(req.Plain)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex plaintext: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexPlaintext, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
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

	logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
	keyHandle, err := pkcs11mgr.FindKeyLabelReturnRandom(req.KeyLabel, *s)
	if err != nil {
		logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(c, ErrorTitleKeyNotFound, ErrorDetailKeyNotExist, ErrorCodeKeyNotFound, http.StatusNotFound, c.Request.URL.Path)
		return
	}

	// Get key attributes
	attr, err := pkcs11mgr.GetObjectAttributes(keyHandle, *s)
	if err != nil {
		logger.AppLog.Errorf("Attributes not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(c, ErrorTitleAttributesNotFound, ErrorDetailAttributesNotFound, ErrorCodeAttributesNotFound, http.StatusNotFound, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Generating initialization vector (IV/nonce) for GCM - 12 bytes recommended")
	iv := make([]byte, 12) // 12 bytes (96 bits) is recommended for GCM
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, ErrorDetailIVGenerationFailed, ErrorCodeIVGenerationFailed, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Encrypting data with AES-GCM")

	// Encrypt with AES-GCM (returns ciphertext + authentication tag)
	ciphertextWithTag, err := pkcs11mgr.EncryptKeyAesGCM(keyHandle, iv, pt, aad, *s)
	if err != nil {
		logger.AppLog.Errorf("AES-GCM encryption failed: %v", err)
		sendProblemDetails(c, ErrorTitleEncryptionFailed, ErrorDetailEncryptionError, ErrorCodeEncryptionError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	safe.Zero(pt) // Clear sensitive data from memory

	logger.AppLog.Info("AES-GCM encryption completed successfully")

	// The output contains ciphertext + 16-byte authentication tag
	// Separate them for the response
	if len(ciphertextWithTag) < 16 {
		logger.AppLog.Errorf("Invalid ciphertext length: %d (expected at least 16 bytes for tag)", len(ciphertextWithTag))
		sendProblemDetails(c, ErrorTitleEncryptionFailed, ErrorDetailInvalidEncryptionOut, ErrorCodeEncryptionError, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	tagLen := 16 // 128-bit tag
	ciphertext := ciphertextWithTag[:len(ciphertextWithTag)-tagLen]
	tag := ciphertextWithTag[len(ciphertextWithTag)-tagLen:]

	ciphertextStr := hex.EncodeToString(ciphertext)
	ivStr := hex.EncodeToString(iv)
	tagStr := hex.EncodeToString(tag)

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Create response using the structured model
	resp := models.EncryptAESGCMResponse{
		Cipher:      ciphertextStr,
		Iv:          ivStr,
		Tag:         tagStr,
		Ok:          true,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
		Id:          attr.Id,
	}

	c.JSON(http.StatusOK, resp)
}
