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
)

// StoreKeyWithPool handles key storage using connection pool
func StoreKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing store key request with pool")

	var req models.StoreKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	label := req.KeyLabel
	id := req.Id
	keyType := req.KeyType

	logger.AppLog.Infof("Decoding key value for label: %s, ID: %s", label, id)
	keyValue, err := hex.DecodeString(req.KeyValue)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode HEX key value: %v", err)
		sendProblemDetails(w, "Bad Request", "The key value in HEX is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.KeyLabel != constants.LABEL_K4_KEY_AES &&
		req.KeyLabel != constants.LABEL_K4_KEY_DES &&
		req.KeyLabel != constants.LABEL_K4_KEY_DES3 {
		logger.AppLog.Errorf("Unsupported key type: %s", req.KeyLabel)
		sendProblemDetails(w, "Bad Request", "The specified key type is not supported", "UNSUPPORTED_KEY_TYPE", http.StatusBadRequest, r.URL.Path)
		return errors.New("unsupported key type")
	}

	var resp models.StoreKeyResponse

	// Use connection pool
	err = pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		logger.AppLog.Infof("Storing key in HSM - Label: %s", label)
		handle, err := mgr.StoreKey(label, keyValue, id, keyType)
		if err != nil {
			logger.AppLog.Errorf("Failed to store key: %v", err)
			sendProblemDetails(w, "Key Storage Failed", "Error storing key in HSM", "KEY_STORAGE_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("Key stored successfully - Handle: %d", handle)

		resp = models.StoreKeyResponse{
			Handle:    int32(handle),
			CipherKey: "", // Initially empty, will be assigned if encryption is possible
		}

		// Try to find the encryption key to encrypt the stored value
		logger.AppLog.Infof("Looking for encryption key: %s", constants.LABEL_ENCRYPTION_KEY)
		findHandle, err := mgr.FindKey(constants.LABEL_ENCRYPTION_KEY, 0)
		if err != nil || findHandle == 0 {
			logger.AppLog.Warnf("Encryption key not found or error: %v. Returning response without encrypted key", err)
			return nil // Continue without encryption
		}

		// Encrypt the stored key value
		logger.AppLog.Info("Encrypting stored key value")
		cipher, err := mgr.EncryptKey(findHandle, nil, keyValue, pkcs11.CKM_AES_CBC_PAD)
		if err != nil {
			logger.AppLog.Warnf("Failed to encrypt stored key value: %v. Returning response without encrypted key", err)
		} else {
			resp.CipherKey = hex.EncodeToString(cipher)
			logger.AppLog.Info("Key value encrypted successfully")
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

// DeleteKeyWithPool handles key deletion using connection pool
func DeleteKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing delete key request with pool")

	var req models.DeleteKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	label := req.KeyLabel
	id := req.Id
	logger.AppLog.Infof("Deleting key with label: %s, ID: %s", label, id)

	var resp models.DeleteKeyResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		// Delete the key from the HSM
		err := mgr.DeleteKey(label, id)
		if err != nil {
			logger.AppLog.Errorf("Failed to delete key: %v", err)
			sendProblemDetails(w, "Key Deletion Failed", "Error deleting key from HSM", "KEY_DELETION_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("Key deleted successfully - Label: %s", label)

		resp = models.DeleteKeyResponse{
			Message:  "Key deleted successfully",
			KeyLabel: label,
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

// UpdateKeyWithPool handles key updates using connection pool
func UpdateKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing update key request with pool")

	var req models.UpdateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	label := req.KeyLabel
	id := req.Id
	keyType := req.KeyType
	logger.AppLog.Infof("Updating key with label: %s, ID: %s, Type: %s", label, id, keyType)

	// Decode the new key value from hexadecimal
	keyValue, err := hex.DecodeString(req.KeyValue)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode HEX key value: %v", err)
		sendProblemDetails(w, "Bad Request", "The key value in HEX is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return err
	}

	var resp models.UpdateKeyResponse

	// Use connection pool
	err = pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		// Update the key in the HSM
		handle, err := mgr.UpdateKey(label, keyValue, req.Id, keyType)
		if err != nil {
			logger.AppLog.Errorf("Failed to update key: %v", err)
			sendProblemDetails(w, "Key Update Failed", "Error updating key in HSM", "KEY_UPDATE_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("Key updated successfully - Label: %s, New Handle: %d", label, handle)

		resp = models.UpdateKeyResponse{
			Message:   "Key updated successfully",
			Handle:    int32(handle),
			KeyLabel:  label,
			CipherKey: "", // Initially empty, will be assigned if encryption is possible
		}

		// Try to find the encryption key to encrypt the new value
		logger.AppLog.Infof("Looking for encryption key: %s", constants.LABEL_ENCRYPTION_KEY)
		findHandle, err := mgr.FindKey(constants.LABEL_ENCRYPTION_KEY, 0)
		if err != nil || findHandle == 0 {
			logger.AppLog.Warnf("Encryption key not found or error: %v. Returning response without encrypted key", err)
			return nil // Continue without encryption
		}

		// Encrypt the new key value
		logger.AppLog.Info("Encrypting updated key value")
		cipher, err := mgr.EncryptKey(findHandle, nil, keyValue, pkcs11.CKM_AES_CBC_PAD)
		if err != nil {
			logger.AppLog.Warnf("Failed to encrypt updated key value: %v. Returning response without encrypted key", err)
		} else {
			resp.CipherKey = hex.EncodeToString(cipher)
			logger.AppLog.Info("Updated key value encrypted successfully")
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
