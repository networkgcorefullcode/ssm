package handlers

import (
	"encoding/json"
	"net/http"

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
func HandleGetDataKeys(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postGetDataKeys(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postGetDataKeys(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing store key request")

	var req models.GetDataKeysRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	label := req.KeyLabel

	logger.AppLog.Infof("Searching key in HSM - using the Label: %s", label)
	handles, err := mgr.FindKeysLabel(label)
	if handles != nil && len(handles) == 0 {
		// Prepare the response
		resp := models.GetDataKeysResponse{
			Keys: make([]models.DataKeyInfo, 0, 0),
		}
		logger.AppLog.Info("Not key found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.AppLog.Errorf("Failed to encode response: %v", err)
			sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
		}
		return
	}
	if err != nil {
		logger.AppLog.Errorf("Failed to search keys: %v", err)
		sendProblemDetails(w, "Key find Failed", "Error searching key in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Keys get successfully")

	objAtr, err := mgr.GetValuesForObjects(handles)
	if err != nil {
		logger.AppLog.Errorf("Failed to get object attributes: %v", err)
		sendProblemDetails(w, "Key get Failed", "Error getting key attributes", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	// Prepare the response
	resp := models.GetDataKeysResponse{
		Keys: make([]models.DataKeyInfo, 0, len(objAtr)),
	}
	for _, attr := range objAtr {
		resp.Keys = append(resp.Keys, models.DataKeyInfo{
			Handle: attr.Handle,
			Id:     attr.Id,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
	}
}
