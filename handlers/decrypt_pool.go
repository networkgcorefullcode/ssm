package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// DecryptWithPool handles decryption using connection pool
func DecryptWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Debugf("Processing decrypt request with pool from %s", r.RemoteAddr)

	var req models.DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Invalid JSON", "Failed to parse request body: "+err.Error(), "bad_json", http.StatusBadRequest, r.URL.Path)
		return err
	}

	// Validate required fields
	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(w, "Validation Error", "Key label is required", "validation_failed", http.StatusBadRequest, r.URL.Path)
		return errors.New("key label is required")
	}

	if req.Cipher == "" {
		logger.AppLog.Error("Ciphertext is required but was empty")
		sendProblemDetails(w, "Validation Error", "Ciphertext is required", "validation_failed", http.StatusBadRequest, r.URL.Path)
		return errors.New("ciphertext is required")
	}

	logger.AppLog.Debugf("Decryption request for key label: %s", req.KeyLabel)

	// Decode ciphertext and IV
	cipher, err := hex.DecodeString(req.Cipher)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode ciphertext hex: %v", err)
		sendProblemDetails(w, "Invalid hex", "Failed to decode ciphertext: "+err.Error(), "bad_hex", http.StatusBadRequest, r.URL.Path)
		return err
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

	var resp models.DecryptResponse

	// Use connection pool
	err = pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		// Find key by label
		keyHandle, err := mgr.FindKey(req.KeyLabel, req.Id)
		if err != nil {
			logger.AppLog.Errorf("Failed to find key by label '%s': %v", req.KeyLabel, err)
			sendProblemDetails(w, "Key Search Failed", "Failed to search for key: "+err.Error(), "key_search_failed", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		var plaintext []byte
		switch req.EncryptionAlgorithm {
		case constants.ALGORITHM_AES128, constants.ALGORITHM_AES256, constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
			plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC_PAD)
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_CBC)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_ECB)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_AES_ECB_ENCRYPT_DATA)
			}
		case constants.ALGORITHM_DES, constants.ALGORITHM_DES_OurUsers:
			plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_CBC_PAD)
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_CBC)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_ECB)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES_ECB_ENCRYPT_DATA)
			}
		case constants.ALGORITHM_DES3, constants.ALGORITHM_DES3_OurUsers:
			plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_CBC_PAD)
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_CBC)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_ECB)
			}
			if err != nil {
				plaintext, err = mgr.DecryptKey(keyHandle, iv, cipher, pkcs11.CKM_DES3_ECB_ENCRYPT_DATA)
			}
		default:
			logger.AppLog.Errorf("Unsupported decryption algorithm: %d", req.EncryptionAlgorithm)
			sendProblemDetails(w, "Bad Request", "Unsupported decryption algorithm", "UNSUPPORTED_ALGORITHM", http.StatusBadRequest, r.URL.Path)
			return errors.New("unsupported algorithm")
		}

		if err != nil {
			logger.AppLog.Errorf("Decryption failed: %v", err)
			sendProblemDetails(w, "Decryption Failed", "Failed to decrypt data: "+err.Error(), "decryption_failed", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("Decryption successful for key '%s', plaintext length: %d bytes", req.KeyLabel, len(plaintext))

		// Prepare response
		resp = models.DecryptResponse{
			Plain: hex.EncodeToString(plaintext),
		}

		// Clear plaintext memory for security
		safe.Zero(plaintext)
		logger.AppLog.Debug("Plaintext memory zeroed for security")

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
		sendProblemDetails(w, "Response Encoding Failed", "Failed to encode response: "+err.Error(), "encoding_failed", http.StatusInternalServerError, r.URL.Path)
		return err
	}

	logger.AppLog.Debug("Decryption response sent successfully")
	return nil
}
