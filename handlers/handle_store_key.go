package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleStoreKey handles key storage requests
// @Summary Store key
// @Description Stores a key in the HSM and optionally encrypts it
// @Tags Key Management
// @Accept json
// @Produce json
// @Param request body models.StoreKeyRequest true "Key data to store"
// @Success 200 {object} models.StoreKeyResponse "Key stored successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /store-key [post]
func HandleStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postStoreKey(mgr, w, r)
	case http.MethodDelete:
		deleteStoreKey(mgr, w, r)
	case http.MethodPut:
		updateStoreKey(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing store key request")

	var req models.StoreKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	label := req.KeyLabel
	id := req.Id
	key_type := req.KeyType
	logger.AppLog.Infof("Decoding key value for label: %s, ID: %s", label, id)
	key_value, err := hex.DecodeString(req.KeyValue)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode HEX key value: %v", err)
		sendProblemDetails(w, "Bad Request", "The key value in HEX is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return
	}

	if req.KeyLabel != constants.LABEL_K4_KEY_AES &&
		req.KeyLabel != constants.LABEL_K4_KEY_DES &&
		req.KeyLabel != constants.LABEL_K4_KEY_DES3 {
		logger.AppLog.Errorf("Unsupported key type: %s", req.KeyLabel)
		sendProblemDetails(w, "Bad Request", "The specified key type is not supported", "UNSUPPORTED_KEY_TYPE", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Storing key in HSM - Label: %s", label)
	handle, err := mgr.StoreKey(label, key_value, []byte(id), key_type)
	if err != nil {
		logger.AppLog.Errorf("Failed to store key: %v", err)
		sendProblemDetails(w, "Key Storage Failed", "Error storing key in HSM", "KEY_STORAGE_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Key stored successfully - Handle: %d", handle)

	resp := models.StoreKeyResponse{
		Handle:    uint(handle),
		CipherKey: nil, // Initially nil, will be assigned if encryption is possible
	}

	// Try to find the encryption key to encrypt the stored value
	logger.AppLog.Infof("Looking for encryption key: %s", constants.LABEL_ENCRYPTION_KEY)
	findHandle, err := mgr.FindKey(constants.LABEL_ENCRYPTION_KEY, "")
	if err != nil || findHandle == 0 {
		logger.AppLog.Warnf("Encryption key not found or error: %v. Returning response without encrypted key", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
		}
		return
	}

	// Encrypt the stored key value
	logger.AppLog.Info("Encrypting stored key value")
	cipher, err := mgr.EncryptKey(findHandle, nil, key_value, pkcs11.CKM_AES_CBC_PAD)
	if err != nil {
		logger.AppLog.Errorf("Failed to encrypt key value: %v. Returning response without encrypted key", err)
		resp.CipherKey = nil
	} else {
		logger.AppLog.Info("Key value encrypted successfully")
		encryptedKeyB64 := hex.EncodeToString(cipher)
		resp.CipherKey = &encryptedKeyB64
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
	}

}

func deleteStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing delete key request")

	var req models.DeleteKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	label := req.KeyLabel
	id := req.Id
	logger.AppLog.Infof("Deleting key with label: %s, ID: %s", label, id)

	// Delete the key from the HSM
	if err := mgr.DeleteKey(label, id); err != nil {
		logger.AppLog.Errorf("Failed to delete key: %v", err)
		sendProblemDetails(w, "Key Deletion Failed", "Error deleting key from HSM", "KEY_DELETION_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Key deleted successfully - Label: %s", label)

	resp := models.DeleteKeyResponse{
		Message:  "Key deleted successfully",
		KeyLabel: label,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
	}
}

func updateStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing update key request")

	var req models.UpdateKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	label := req.KeyLabel
	id := req.Id
	keyType := req.KeyType
	logger.AppLog.Infof("Updating key with label: %s, ID: %s, Type: %s", label, id, keyType)

	// Decode the new key value from hexadecimal
	keyValue, err := hex.DecodeString(req.KeyValue)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex key value: %v", err)
		sendProblemDetails(w, "Bad Request", "The key value in hexadecimal is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return
	}

	// Update the key in the HSM
	handle, err := mgr.UpdateKey(label, keyValue, []byte(id), keyType)
	if err != nil {
		logger.AppLog.Errorf("Failed to update key: %v", err)
		sendProblemDetails(w, "Key Update Failed", "Error updating key in HSM", "KEY_UPDATE_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Key updated successfully - Label: %s, New Handle: %d", label, handle)

	resp := models.UpdateKeyResponse{
		Message:   "Key updated successfully",
		Handle:    uint(handle),
		KeyLabel:  label,
		CipherKey: nil, // Initially nil, will be assigned if encryption is possible
	}

	// Try to find the encryption key to encrypt the new value
	logger.AppLog.Infof("Looking for encryption key: %s", constants.LABEL_ENCRYPTION_KEY)
	findHandle, err := mgr.FindKey(constants.LABEL_ENCRYPTION_KEY, "")
	if err != nil || findHandle == 0 {
		logger.AppLog.Warnf("Encryption key not found or error: %v. Returning response without encrypted key", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
		}
		return
	}

	// Encrypt the new key value
	logger.AppLog.Info("Encrypting updated key value")
	cipher, err := mgr.EncryptKey(findHandle, nil, keyValue, pkcs11.CKM_AES_CBC_PAD)
	if err != nil {
		logger.AppLog.Errorf("Failed to encrypt updated key value: %v. Returning response without encrypted key", err)
		resp.CipherKey = nil
	} else {
		logger.AppLog.Info("Updated key value encrypted successfully")
		encryptedKeyB64 := hex.EncodeToString(cipher)
		resp.CipherKey = &encryptedKeyB64
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
	}
}
