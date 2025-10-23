package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// EncryptWithPool handles encryption using connection pool
func EncryptWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing encrypt request")

	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	logger.AppLog.Infof("Decoding hex plaintext for key: %s", req.KeyLabel)
	pt, err := hex.DecodeString(req.Plain)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex plaintext: %v", err)
		sendProblemDetails(w, "Bad Request", "The hex data is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return err
	}
	defer safe.Zero(pt) // Ensure sensitive data is cleared

	var resp models.EncryptResponse

	// Use connection pool
	err = pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
		keyHandle, err := mgr.FindKeyLabelReturnRandom(req.KeyLabel)
		if err != nil {
			logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
			sendProblemDetails(w, "Key Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
			return err
		}

		atrr, err := mgr.GetObjectAttributes(keyHandle)
		if err != nil {
			logger.AppLog.Errorf("Attributes not found: %s, error: %v", req.KeyLabel, err)
			sendProblemDetails(w, "Attributes Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
			return err
		}

		logger.AppLog.Info("Generating initialization vector (IV)")
		var size int
		if req.EncryptionAlgorithm == constants.ALGORITHM_DES3_OurUsers || req.EncryptionAlgorithm == constants.ALGORITHM_DES_OurUsers {
			size = 8
		} else if req.EncryptionAlgorithm == constants.ALGORITHM_AES128_OurUsers || req.EncryptionAlgorithm == constants.ALGORITHM_AES256_OurUsers {
			size = 16
		}

		iv := make([]byte, size)
		if err := safe.RandRead(iv); err != nil {
			logger.AppLog.Errorf("Failed to generate IV: %v", err)
			sendProblemDetails(w, "Internal Server Error", "Error generating initialization vector", "IV_GENERATION_FAILED", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Info("Encrypting data")

		var ciphertext []byte
		switch req.EncryptionAlgorithm {
		case constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC_PAD)
			if err != nil {
				logger.AppLog.Errorf("Encryption with PAD failed: %v", err)
				ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC)
				if err != nil {
					logger.AppLog.Errorf("Encryption failed: %v", err)
					sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
					return err
				}
			}
		case constants.ALGORITHM_DES3_OurUsers:
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC_PAD)
			if err != nil {
				logger.AppLog.Errorf("Encryption with PAD failed: %v", err)
				ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC)
				if err != nil {
					logger.AppLog.Errorf("Encryption failed: %v", err)
					sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
					return err
				}
			}
		case constants.ALGORITHM_DES_OurUsers:
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC_PAD)
			if err != nil {
				logger.AppLog.Errorf("Encryption with PAD failed: %v", err)
				ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC)
				if err != nil {
					logger.AppLog.Errorf("Encryption failed: %v", err)
					sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
					return err
				}
			}
		default:
			logger.AppLog.Errorf("Unsupported encryption algorithm: %d", req.EncryptionAlgorithm)
			sendProblemDetails(w, "Bad Request", "Unsupported encryption algorithm", "UNSUPPORTED_ALGORITHM", http.StatusBadRequest, r.URL.Path)
			return errors.New("unsupported algorithm")
		}

		logger.AppLog.Info("Encryption completed successfully")

		// Create response
		timeCreated := time.Now()
		resp = models.EncryptResponse{
			Cipher:      hex.EncodeToString(ciphertext),
			Iv:          hex.EncodeToString(iv),
			Ok:          true,
			TimeCreated: timeCreated,
			TimeUpdated: timeCreated,
			Id:          atrr.Id,
		}

		return nil
	})

	if err != nil {
		// Error was already handled inside the WithConnection function
		return err
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		return err
	}

	return nil
}
