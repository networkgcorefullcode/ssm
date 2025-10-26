package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func HandleGetAllKeys(c *gin.Context) {
	logger.AppLog.Info("Processing get all keys request")
	//// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	// Find all keys grouped by label
	logger.AppLog.Info("Searching all keys in HSM")
	keysByLabel, err := pkcs11mgr.FindAllKeys(*s)
	if err != nil && err.Error() == "error Key With The Label Not Found" {
		// Prepare the response
		resp := models.GetAllKeysResponse{}
		logger.AppLog.Info("Not key found")
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		logger.AppLog.Errorf("Failed to search all keys: %v", err)
		sendProblemDetails(c, "Key Search Failed", "Error searching all keys in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
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

		objAttrs, err := pkcs11mgr.GetValuesForObjects(handles, *s)
		if err != nil {
			logger.AppLog.Errorf("Failed to get object attributes for label %s: %v", label, err)
			sendProblemDetails(c, "Key Attributes Failed", "Error getting key attributes", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
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

	c.JSON(http.StatusOK, resp)
}
