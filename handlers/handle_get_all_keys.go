package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleGetAllKeys handles requests to get all keys from HSM
// @Summary Get all keys
// @Description Retrieves all keys from the HSM grouped by label
// @Tags Key Management
// @Accept json
// @Produce json
// @Success 200 {object} models.GetAllKeysResponse "All keys retrieved successfully"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /get-all-keys [post]
func HandleGetAllKeys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postGetAllKeys(w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postGetAllKeys(w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing get all keys request")
	// init the pkcs manager
	mgr, err := pkcs11mgr.New(factory.SsmConfig.Configuration.PkcsPath,
		uint(factory.SsmConfig.Configuration.LotsNumber),
		factory.SsmConfig.Configuration.Pin)
	if err != nil {
		logger.AppLog.Errorf("Failed to create PKCS11 manager: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Failed to initialize PKCS11 manager", "PKCS_INIT_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	err = mgr.OpenSession()
	if err != nil {
		logger.AppLog.Errorf("Failed to OpenSession PKCS11 manager: %v", err)
		sendProblemDetails(w, "Internal Server Error", "The pkcs session have a error during stablishment", "PKCS_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	defer mgr.CloseSession()
	defer mgr.Finalize()

	// Find all keys grouped by label
	logger.AppLog.Info("Searching all keys in HSM")
	keysByLabel, err := mgr.FindAllKeys()
	if err != nil && err.Error() == "Key with the label not found" {
		// Prepare the response
		// Prepare the response
		// Prepare the response
		resp := models.GetAllKeysResponse{}
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
		logger.AppLog.Errorf("Failed to search all keys: %v", err)
		sendProblemDetails(w, "Key Search Failed", "Error searching all keys in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Found keys in %d labels", len(keysByLabel))

	// Prepare the response
	resp := models.GetAllKeysResponse{
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
			return
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

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, r.URL.Path)
	}
}
