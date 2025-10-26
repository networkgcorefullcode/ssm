package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleGetDataKey
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
func HandleGetDataKey(c *gin.Context) {
	logger.AppLog.Info("Processing get data key request")
	// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	var req models.GetKeyRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	label := req.KeyLabel

	logger.AppLog.Infof("Searching key in HSM - using the Label: %s", label)
	handle, err := pkcs11mgr.FindKey(label, req.Id, *s)
	if err != nil && err.Error() == "error Key With The Label Not Found" {
		resp := models.GetKeyResponse{}
		logger.AppLog.Info("Not key found")
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		logger.AppLog.Errorf("Failed to search keys: %v", err)
		sendProblemDetails(c, "Key find Failed", "Error searching key in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Key get successfully")

	objAtr, err := pkcs11mgr.GetObjectAttributes(handle, *s)
	if err != nil {
		logger.AppLog.Errorf("Failed to get object attribute: %v", err)
		sendProblemDetails(c, "Key get Failed", "Error getting key attribute", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	resp := models.GetKeyResponse{
		KeyInfo: models.DataKeyInfo{
			Handle: objAtr.Handle,
			Id:     objAtr.Id,
		},
	}

	c.JSON(http.StatusOK, resp)
}
