package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// GetAllKeysWithPool handles requests to get all keys from HSM using connection pool
func GetAllKeysWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing get all keys request with pool")

	var resp models.GetAllKeysResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		// Find all keys grouped by label
		logger.AppLog.Info("Searching all keys in HSM")
		keysByLabel, err := mgr.FindAllKeys()
		if err != nil && err.Error() == "Key with the label not found" {
			// Prepare empty response
			resp = models.GetAllKeysResponse{}
			logger.AppLog.Info("No keys found")
			return nil
		}
		if err != nil {
			logger.AppLog.Errorf("Failed to search all keys: %v", err)
			sendProblemDetails(w, "Key Search Failed", "Error searching all keys in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("Found keys in %d labels", len(keysByLabel))

		// Prepare the response
		resp = models.GetAllKeysResponse{
			KeysByLabel: make(map[string][]models.DataKeyInfo),
			TotalKeys:   0,
			TotalLabels: int32(len(keysByLabel)),
		}

		// Process each label and its keys
		for label, handles := range keysByLabel {
			logger.AppLog.Infof("Processing label: %s with %d keys", label, len(handles))

			objAttrs, err := mgr.GetValuesForObjects(handles)
			if err != nil {
				logger.AppLog.Errorf("Failed to get object attributes for label %s: %v", label, err)
				sendProblemDetails(w, "Key Attributes Failed", "Error getting key attributes", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
				return err
			}

			// Convert to DataKeyInfo
			keysInfo := make([]models.DataKeyInfo, 0, len(objAttrs))
			for _, attr := range objAttrs {
				keysInfo = append(keysInfo, models.DataKeyInfo{
					Handle: attr.Handle,
					Id:     attr.Id,
				})
			}

			resp.KeysByLabel[label] = keysInfo
			resp.TotalKeys += int32(len(keysInfo))
		}

		logger.AppLog.Infof("Successfully retrieved %d keys across %d labels", resp.TotalKeys, resp.TotalLabels)
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
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
		return err
	}

	return nil
}

// GetDataKeyWithPool handles requests to get a specific key using connection pool
func GetDataKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing get data key request with pool")

	var req models.GetKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(w, "Bad Request", "Key label is required", "MISSING_KEY_LABEL", http.StatusBadRequest, r.URL.Path)
		return errors.New("key label is required")
	}

	label := req.KeyLabel
	var resp models.GetKeyResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		logger.AppLog.Infof("Searching key in HSM - using the Label: %s", label)
		handle, err := mgr.FindKey(label, req.Id)
		if err != nil && err.Error() == "Key with the label not found" {
			// Prepare empty response
			resp = models.GetKeyResponse{}
			logger.AppLog.Info("No key found")
			return nil
		}
		if err != nil {
			logger.AppLog.Errorf("Failed to search keys: %v", err)
			sendProblemDetails(w, "Key find Failed", "Error searching key in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Info("Key found successfully")

		objAtr, err := mgr.GetObjectAttributes(handle)
		if err != nil {
			logger.AppLog.Errorf("Failed to get object attribute: %v", err)
			sendProblemDetails(w, "Key get Failed", "Error getting key attribute", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		// Prepare the response
		resp = models.GetKeyResponse{
			KeyInfo: models.DataKeyInfo{
				Handle: objAtr.Handle,
				Id:     objAtr.Id,
			},
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
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
		return err
	}

	return nil
}

// GetDataKeysWithPool handles requests to get multiple keys by label using connection pool
func GetDataKeysWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing get data keys request with pool")

	var req models.GetDataKeysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.KeyLabel == "" {
		logger.AppLog.Error("Key label is required but was empty")
		sendProblemDetails(w, "Bad Request", "Key label is required", "MISSING_KEY_LABEL", http.StatusBadRequest, r.URL.Path)
		return errors.New("key label is required")
	}

	label := req.KeyLabel
	var resp models.GetDataKeysResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		logger.AppLog.Infof("Searching keys in HSM - using the Label: %s", label)
		handles, err := mgr.FindKeysLabel(label)
		if err != nil && err.Error() == "Key with the label not found" {
			// Prepare empty response
			resp = models.GetDataKeysResponse{
				Keys: make([]models.DataKeyInfo, 0, 0),
			}
			logger.AppLog.Info("No keys found")
			return nil
		}
		if err != nil {
			logger.AppLog.Errorf("Failed to search keys: %v", err)
			sendProblemDetails(w, "Key find Failed", "Error searching keys in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Info("Keys found successfully")

		objAtr, err := mgr.GetValuesForObjects(handles)
		if err != nil {
			logger.AppLog.Errorf("Failed to get object attributes: %v", err)
			sendProblemDetails(w, "Key get Failed", "Error getting key attributes", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		// Prepare the response
		resp = models.GetDataKeysResponse{
			Keys: make([]models.DataKeyInfo, 0, len(objAtr)),
		}
		for _, attr := range objAtr {
			resp.Keys = append(resp.Keys, models.DataKeyInfo{
				Handle: attr.Handle,
				Id:     attr.Id,
			})
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
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
		return err
	}

	return nil
}
