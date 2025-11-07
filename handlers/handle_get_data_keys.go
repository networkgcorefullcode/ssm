package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleGetDataKeys
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
func HandleGetDataKeys(c *gin.Context) {
	logger.AppLog.Info("Processing store key request")
	// init the session
	s := mgr.GetSession()
	defer mgr.LogoutSession(s)

	var req models.GetDataKeysRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(c, ErrorTitleBadRequest, ErrorDetailInvalidJSON, ErrorCodeInvalidJSON, http.StatusBadRequest, c.Request.URL.Path)
		return
	}

	label := req.KeyLabel

	logger.AppLog.Infof("Searching key in HSM - using the Label: %s", label)
	handles, err := pkcs11mgr.FindKeysLabel(label, *s)
	if err != nil && err.Error() == constants.ERROR_STRING_KEY_NOT_FOUND {
		resp := models.GetDataKeysResponse{
			Keys: make([]models.DataKeyInfo, 0),
		}
		logger.AppLog.Info("Not key found")
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		logger.AppLog.Errorf("Failed to search keys: %v", err)
		sendProblemDetails(c, "Key find Failed", "Error searching key in HSM", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	logger.AppLog.Info("Keys get successfully")

	objAtr, err := pkcs11mgr.GetValuesForObjects(handles, *s)
	if err != nil {
		logger.AppLog.Errorf("Failed to get object attributes: %v", err)
		sendProblemDetails(c, "Key get Failed", "Error getting key attributes", "KEY_GET_ERROR", http.StatusInternalServerError, c.Request.URL.Path)
		return
	}

	resp := models.GetDataKeysResponse{
		Keys: make([]models.DataKeyInfo, 0, len(objAtr)),
	}
	for _, attr := range objAtr {
		resp.Keys = append(resp.Keys, models.DataKeyInfo{
			Handle: attr.Handle,
			Id:     attr.Id,
		})
	}

	c.JSON(http.StatusOK, resp)
}
