package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/database"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// AuditLog represents a structured audit log entry
type AuditLog struct {
	Start time.Time `json:"start_time"`
	// UserID     string    `json:"user_id,omitempty"`
	Action     string `json:"action"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent,omitempty"`
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id,omitempty"`
	Duration   int64  `json:"duration_ms,omitempty"` // Duration in milliseconds
	Error      string `json:"error,omitempty"`
	Signature  string `json:"signature,omitempty"`
}

// Map common patterns to actions
var ActionMap map[string]string = map[string]string{
	"POST /crypto/encrypt":           "ENCRYPT_DATA",
	"POST /crypto/decrypt":           "DECRYPT_DATA",
	"POST /crypto/generate-aes-key":  "GENERATE_AES_KEY",
	"POST /crypto/generate-des-key":  "GENERATE_DES_KEY",
	"POST /crypto/generate-des3-key": "GENERATE_DES3_KEY",
	"POST /crypto/store-key":         "STORE_KEY",
	"PUT /crypto/store-key":          "UPDATE_KEY",
	"DELETE /crypto/store-key":       "DELETE_KEY",
	"POST /crypto/get-data-key":      "GET_KEY",
	"POST /crypto/get-data-keys":     "GET_KEYS",
	"POST /crypto/get-all-keys":      "GET_ALL_KEYS",
	"POST /crypto/health-check":      "HEALTH_CHECK",
	"POST /login":                    "USER_LOGIN",
}

func AuditRequest(c *gin.Context) {
	// Capture data before processing the request
	start := time.Now()
	c.Next() // Process the request

	// Calculate duration
	duration := time.Since(start)
	// Get or generate Request ID
	requestID := c.GetString("request_id")
	if requestID == "" {
		requestID = generateContextualRequestID()
		c.Set("request_id", requestID)
	}
	// Log data
	logEntry := AuditLog{
		Start: start,
		// UserID:     c.GetString("user_id"), // Assume user_id is set in context
		Action:     determineAction(c), // E.g.: based on route+method
		Method:     c.Request.Method,
		Path:       c.Request.URL.Path,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		StatusCode: c.Writer.Status(),
		RequestID:  c.GetString("request_id"),
		Duration:   duration.Milliseconds(),
	}

	// Capture errors if they exist
	if len(c.Errors) > 0 {
		logEntry.Error = c.Errors.String()
	}

	// TODO: Implement audit log persistence (database, file, etc.)
	// For now, just log to application logger
	logAuditEntry(logEntry)
}

func generateContextualRequestID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)

	return fmt.Sprintf("ssm-%s-%d-%s",
		hostname[:min(len(hostname), 8)],
		time.Now().Unix(),
		hex.EncodeToString(bytes))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// logAuditEntry logs the audit entry
// TODO: Implement persistence to database or file system
func logAuditEntry(entry AuditLog) {
	//use the application logger
	logInLogger(entry)

	entry.Signature = "" // Clear signature before signing

	// getting pkcs11 manager and session from SsmServer
	session := mgr.GetSession()
	auditPrivateKey := pkcs11mgr.GetAuditPrivateKey()

	var err error
	entry.Signature, err = signAuditLog(entry, &session.Handle, session.Ctx, auditPrivateKey)
	if err != nil {
		logger.AppLog.Errorf("Failed to sign audit log: %v", err)
		return
	}

	go database.InsertData(database.Client, factory.SsmConfig.Configuration.Mongodb.DBName, database.CollAuditLogs, entry)
}

func signAuditLog(logData any, session *pkcs11.SessionHandle, mgr *pkcs11.Ctx, privKey pkcs11.ObjectHandle) (string, error) {
	data, err := json.Marshal(logData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)

	err = mgr.SignInit(*session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)}, privKey)
	if err != nil {
		return "", err
	}

	signature, err := mgr.Sign(*session, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// determineAction determines the action based on the request path and method
func determineAction(c *gin.Context) string {
	path := c.Request.URL.Path
	method := c.Request.Method

	key := method + " " + path
	if action, exists := ActionMap[key]; exists {
		return action
	}

	// Default action
	return method + "_" + path
}

// logInLogger logs the audit entry to the application logger
func logInLogger(entry AuditLog) {
	if entry.StatusCode >= 400 {
		logger.AppLog.Warnf("[AUDIT] %s | %s %s | IP: %s | Status: %d | Duration: %dms | User: %s | RequestID: %s | Error: %s",
			entry.Start.Format(time.RFC3339),
			entry.Method,
			entry.Path,
			entry.IP,
			entry.StatusCode,
			entry.Duration,
			// entry.UserID,
			entry.RequestID,
			entry.Error,
		)
	} else {
		logger.AppLog.Infof("[AUDIT] %s | %s %s | IP: %s | Status: %d | Duration: %dms | User: %s | RequestID: %s",
			entry.Start.Format(time.RFC3339),
			entry.Method,
			entry.Path,
			entry.IP,
			entry.StatusCode,
			entry.Duration,
			// entry.UserID,
			entry.RequestID,
		)
	}
}
