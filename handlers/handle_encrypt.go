package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// HandleEncrypt handles encryption requests
// @Summary Encrypt data
// @Description Encrypts data using an AES key stored in the HSM
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.EncryptRequest true "Data to encrypt"
// @Success 200 {object} models.EncryptResponse "Data encrypted successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 404 {object} models.ProblemDetails "Key not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /encrypt [post]
func HandleEncrypt(c *gin.Context) {
	// init the session
	s := mgr.GetSession()
	//

	defer mgr.LogoutSession(s)

	logger.AppLog.Info("Processing encrypt request")

	var req models.EncryptRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidJSON, ErrorCodeInvalidJSON, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("Decoding hex plaintext for key: %s", req.KeyLabel)
	pt, err := hex.DecodeString(req.Plain)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex plaintext: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidHexData, ErrorCodeInvalidHex, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
	keyHandle, err := pkcs11mgr.FindKeyLabelReturnRandom(req.KeyLabel, *s)
	if err != nil {
		logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(c, ErrorTitleKeyNotFound, ErrorDetailKeyNotExist, ErrorCodeKeyNotFound, http.StatusNotFound, c.Request.URL.Path)
		return
	}
	atrr, err := pkcs11mgr.GetObjectAttributes(keyHandle, *s)
	if err != nil {
		logger.AppLog.Errorf("Atributes not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(c, ErrorTitleAttributesNotFound, ErrorDetailAttributesNotFound, ErrorCodeAttributesNotFound, http.StatusNotFound, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Generating initialization vector (IV)")
	var size int
	switch req.EncryptionAlgorithm {
	case constants.ALGORITHM_DES3_OurUsers, constants.ALGORITHM_DES_OurUsers:
		size = 8
	case constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
		size = 16
	}
	iv := make([]byte, size)
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		sendProblemDetails(c, ErrorTitleInternalServerError, ErrorDetailIVGenerationFailed, ErrorCodeIVGenerationFailed, http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Encrypting data")

	var ciphertext []byte
	switch req.EncryptionAlgorithm {
	case constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
		ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC_PAD, *s)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC, *s)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(c, ErrorTitleEncryptionFailed, ErrorDetailEncryptionError, ErrorCodeEncryptionError, http.StatusInternalServerError, c.Request.URL.Path)
				return
			}
		}
	case constants.ALGORITHM_DES3_OurUsers:
		ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC_PAD, *s)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC, *s)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(c, ErrorTitleEncryptionFailed, ErrorDetailEncryptionError, ErrorCodeEncryptionError, http.StatusInternalServerError, c.Request.URL.Path)
				return
			}
		}
	case constants.ALGORITHM_DES_OurUsers:
		ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC_PAD, *s)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = pkcs11mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC, *s)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(c, ErrorTitleEncryptionFailed, ErrorDetailEncryptionError, ErrorCodeEncryptionError, http.StatusInternalServerError, c.Request.URL.Path)
				return
			}
		}
	default:
		logger.AppLog.Errorf("Unsupported encryption algorithm: %s", req.EncryptionAlgorithm)
		sendProblemDetails(c, ErrorTitleBadRequest, "The specified encryption algorithm is not supported", "UNSUPPORTED_ALGORITHM", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	safe.Zero(pt) // Clear sensitive data from memory

	logger.AppLog.Info("Encryption completed successfully")

	ciphertextStr := hex.EncodeToString(ciphertext)
	ivStr := hex.EncodeToString(iv)
	ok := true

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Create response using the structured model
	resp := models.EncryptResponse{
		Cipher:      ciphertextStr,
		Iv:          ivStr,
		Ok:          ok,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
		Id:          atrr.Id,
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// if err := json.NewEncoder(w).Encode(resp); err != nil {
	// 	logger.AppLog.Errorf("Failed to encode response: %v", err)
	// }

	c.JSON(http.StatusCreated, resp)
}
