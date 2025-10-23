package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandlePoolStats returns PKCS11 connection pool statistics
// @Summary Get PKCS11 connection pool statistics
// @Description Returns statistics about the PKCS11 connection pool usage
// @Tags Monitoring
// @Produce json
// @Success 200 {object} pkcs11mgr.PoolStats "Pool statistics"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /pool/stats [get]
func HandlePoolStats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPoolStats(w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func getPoolStats(w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Getting PKCS11 pool statistics")

	stats := pkcs11mgr.GetPoolStats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logger.AppLog.Errorf("Failed to encode pool stats response: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Failed to encode response", "ENCODING_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Pool stats returned: %+v", stats)
}
