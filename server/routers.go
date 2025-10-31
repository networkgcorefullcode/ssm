package server

import (
	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/handlers"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/server/middleware"
)

// CreateGinRouter sets up routes using Gin and wraps existing net/http handlers for compatibility.
func CreateGinRouter() *gin.Engine {
	// Use ReleaseMode unless verbose debugging is required; Gin still logs via its middleware.
	// Mode can be adjusted by env GIN_MODE if needed.
	r := gin.New()
	r.Use(gin.Recovery()) // recover from panics and write 500
	r.Use(gin.Logger())   // basic logging middleware

	// create router crypto group
	rc := r.Group("/crypto")

	// middlewares for security, logging, tracing, etc.
	if factory.SsmConfig.Configuration.IsSecure {
		logger.AppLog.Info("Configuring secure middlewares")
		r.Use(middleware.AuditRequest)
		middleware.ConfigureCORS(r) // configure CORS if needed
		// r.Use(middleware.ValidateRequest)     // validate request schema, headers, etc.
		r.Use(middleware.SecureRequest)          // secure middleware for headers, rate limiting, etc.
		rc.Use(middleware.AuthenticateRequest()) // authentication middleware
	}

	// Endpoints

	r.POST("/login", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /login request")
		handlers.HandleLogin(c)
	})

	// HealthCheck endpoint (GET recommended)
	rc.GET("/health-check", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /health-check request")
		handlers.HandleHealthCheck(c)
	})

	// Encrypt endpoints POST
	rc.POST("/encrypt", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /encrypt request")
		handlers.HandleEncrypt(c)
	})

	// Decrypt endpoints POST
	rc.POST("/decrypt", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /decrypt request")
		handlers.HandleDecrypt(c)
	})

	// Store Key endpoints POST
	rc.POST("/store-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /store-key request")
		handlers.HandleStoreKey(c)
	})
	rc.PUT("/store-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /store-key request")
		handlers.HandleStoreKey(c)
	})
	rc.DELETE("/store-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /store-key request")
		handlers.HandleStoreKey(c)
	})

	// Generate Key endpoints POST
	rc.POST("/generate-aes-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /generate-aes-key request")
		handlers.HandleGenerateAESKey(c)
	})

	rc.POST("/generate-des3-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /generate-des3-key request")
		handlers.HandleGenerateDES3Key(c)
	})

	rc.POST("/generate-des-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /generate-des-key request")
		handlers.HandleGenerateDESKey(c)
	})

	// Synchronization handlers
	rc.POST("/get-data-keys", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /get-data-keys request")
		handlers.HandleGetDataKeys(c)
	})

	rc.POST("/get-key", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /get-keys request")
		handlers.HandleGetDataKey(c)
	})

	rc.POST("/get-all-keys", func(c *gin.Context) {
		logger.AppLog.Debugf("Received /get-all-keys request")
		handlers.HandleGetAllKeys(c)
	})

	return r
}
