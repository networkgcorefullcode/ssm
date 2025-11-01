package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// @title        Decrypt Data API
// @version 	 1.0.0
// @description  Decrypt a Key using simetrics ALGORITHMs as AES 128 and AES 256, DES, 3DES
// @Accept       json
// @Produce      json
// @Param        request  body      models.DecryptRequest  true  "Data to decrypt"
// @Success      200      {object}  models.DecryptResponse "Data decrypted successfully"
// @Failure      400      {object}  models.ProblemDetails  "Validation error or invalid JSON"
// @Failure      404      {object}  models.ProblemDetails  "Key not found"
// @Failure      405      {object}  models.ProblemDetails  "HTTP method not allowed"
// @Failure      500      {object}  models.ProblemDetails  "Internal server error"
// @Router       /decrypt [post]
func HandleDecrypt(c *gin.Context) {
	logger.AppLog.Debugf("Received decrypt request from %s", c.ClientIP())

	logger.AppLog.Debugf("Processing decrypt request for %s", c.Request.URL.Path)
	// init the session
	s := mgr.GetSession()
	//

	defer mgr.LogoutSession(s)

	var req models.DecryptRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, "Invalid JSON", "Failed to parse request body: "+err.Error(), "bad_json", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	logger.AppLog.Debugf("Decryption request for key label: %s", req.KeyLabel)

	// Validate required fields
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(c, "Validation Error", "Key label is required", "validation_failed", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if req.Cipher == "" {
		logger.AppLog.Error("Ciphertext is required but was empty")
		sendProblemDetails(c, "Validation Error", "Ciphertext is required", "validation_failed", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	// Decode ciphertext and IV
	cipher, err := hex.DecodeString(req.Cipher)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode ciphertext hex: %v", err)
		sendProblemDetails(c, "Invalid hex", "Failed to decode ciphertext: "+err.Error(), "bad_hex", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	iv, err := hex.DecodeString(req.Iv)

	if err != nil {
		logger.AppLog.Errorf("Failed to decode IV hex: %v", err)
		iv = nil
	}
	if req.Iv == "" {
		logger.AppLog.Info("iv is empty")
		iv = nil
	}

	logger.AppLog.Debugf("Decoded ciphertext length: %d bytes, IV length: %d bytes", len(cipher), len(iv))

	// Find key by label
	keyHandle, err := pkcs11mgr.FindKey(req.KeyLabel, req.Id, *s)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key by label '%s': %v", req.KeyLabel, err)
		sendProblemDetails(c, "Key Search Failed", "Failed to search for key: "+err.Error(), "key_search_failed", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	var plaintext []byte
	switch req.EncryptionAlgorithm {
	case constants.ALGORITHM_AES128, constants.ALGORITHM_AES256, constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
		if len(iv) == 16 {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC_PAD, *s)
			if err != nil {
				plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC, *s)
			}
		} else if iv == nil {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_ECB, *s)
		}
		// if err != nil {
		// 	plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_ECB_ENCRYPT_DATA)
		// }
	case constants.ALGORITHM_DES, constants.ALGORITHM_DES_OurUsers:
		if len(iv) == 8 {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_CBC_PAD, *s)
			if err != nil {
				plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_CBC, *s)
			}
		} else if iv == nil {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_ECB, *s)
		}
		// if err != nil {
		// 	plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_ECB_ENCRYPT_DATA)
		// }
	case constants.ALGORITHM_DES3, constants.ALGORITHM_DES3_OurUsers:
		if len(iv) == 8 {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_CBC_PAD, *s)
			if err != nil {
				plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_CBC, *s)
			}
		} else if iv == nil {
			plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_ECB, *s)
		}
		// if err != nil {
		// 	plaintext, err = pkcs11mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_ECB_ENCRYPT_DATA)
		// }
	default:
		logger.AppLog.Errorf("Unsupported decryption algorithm: %d", req.EncryptionAlgorithm)
		sendProblemDetails(c, "Bad Request", "Unsupported decryption algorithm", "UNSUPPORTED_ALGORITHM", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	if err != nil {
		logger.AppLog.Errorf("Decryption failed: %v", err)
		sendProblemDetails(c, "Decryption Failed", "Failed to decrypt data: "+err.Error(), "decryption_failed", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}
	if len(plaintext) == 0 {
		logger.AppLog.Error("Decryption resulted in empty plaintext")
		sendProblemDetails(c, "Decryption Failed", "Decryption resulted in empty plaintext", "empty_plaintext", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Infof("Decryption successful for key '%s', plaintext length: %d bytes", req.KeyLabel, len(plaintext))

	// Prepare response
	c.Header("Content-Type", "application/json")
	resp := models.DecryptResponse{
		Plain: hex.EncodeToString(plaintext),
	}

	// Clear plaintext memory for security
	safe.Zero(plaintext)
	logger.AppLog.Debug("Plaintext memory zeroed for security")

	c.JSON(http.StatusOK, resp)

	logger.AppLog.Debug("Decryption response sent successfully")
}
