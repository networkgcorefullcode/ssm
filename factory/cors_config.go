package factory

import (
	"time"

	"github.com/networkgcorefullcode/ssm/logger"
)

// CORS configuration based on gin-contrib/cors
type CORS struct {
	// AllowAllOrigins allows all origins. Cannot be used with AllowOrigins, AllowOriginFunc, or AllowOriginWithContextFunc
	AllowAllOrigins *bool `yaml:"allowAllOrigins,omitempty"`

	// AllowOrigins is a list of allowed origins. Supports exact match, *, and wildcards
	AllowOrigins []string `yaml:"allowOrigins,omitempty"`

	// AllowMethods is a list of allowed HTTP methods
	// Default: ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"]
	AllowMethods []string `yaml:"allowMethods,omitempty"`

	// AllowPrivateNetwork adds Private Network Access CORS header
	AllowPrivateNetwork *bool `yaml:"allowPrivateNetwork,omitempty"`

	// AllowHeaders is a list of non-simple headers permitted in requests
	AllowHeaders []string `yaml:"allowHeaders,omitempty"`

	// AllowCredentials allows cookies, HTTP auth, or client certs. Only if precise origins are used
	AllowCredentials *bool `yaml:"allowCredentials,omitempty"`

	// ExposeHeaders is a list of headers exposed to the browser
	ExposeHeaders []string `yaml:"exposeHeaders,omitempty"`

	// MaxAge cache time for preflight requests in seconds
	// Will be converted to time.Duration when used
	MaxAge *int `yaml:"maxAge,omitempty"`

	// AllowWildcard enables wildcards in origins (e.g. https://*.example.com)
	AllowWildcard *bool `yaml:"allowWildcard,omitempty"`

	// AllowBrowserExtensions allows browser extension schemes as origins (e.g. chrome-extension://)
	AllowBrowserExtensions *bool `yaml:"allowBrowserExtensions,omitempty"`

	// CustomSchemas are additional allowed URI schemes (e.g. tauri://)
	CustomSchemas []string `yaml:"customSchemas,omitempty"`

	// AllowWebSockets allows ws:// and wss:// schemas
	AllowWebSockets *bool `yaml:"allowWebSockets,omitempty"`

	// AllowFiles allows file:// origins (dangerous; use only if necessary)
	AllowFiles *bool `yaml:"allowFiles,omitempty"`

	// OptionsResponseStatusCode custom status code for OPTIONS responses
	// Default: 204
	OptionsResponseStatusCode *int `yaml:"optionsResponseStatusCode,omitempty"`
}

// GetMaxAgeDuration converts MaxAge from seconds to time.Duration
func (c *CORS) GetMaxAgeDuration() time.Duration {
	if c.MaxAge == nil {
		return 12 * time.Hour // Default value
	}
	return time.Duration(*c.MaxAge) * time.Second
}

// GetOptionsStatusCode returns the OPTIONS status code or default 204
func (c *CORS) GetOptionsStatusCode() int {
	if c.OptionsResponseStatusCode == nil {
		return 204
	}
	return *c.OptionsResponseStatusCode
}

// initializeCORSConfig initializes CORS configuration with secure defaults
func initializeCORSConfig() {
	// If CORS is nil, create it with disabled state
	if SsmConfig.Configuration.CORS == nil {
		enabled := false
		SsmConfig.Configuration.CORS = &CORS{
			AllowAllOrigins: &enabled,
		}
		logger.CfgLog.Info("CORS configuration not set. CORS will be disabled")
		return
	}

	// If AllowAllOrigins is not set, default to false for security
	if SsmConfig.Configuration.CORS.AllowAllOrigins == nil {
		disabled := false
		SsmConfig.Configuration.CORS.AllowAllOrigins = &disabled
		logger.CfgLog.Info("CORS AllowAllOrigins not set. Defaulting to false for security")
	}

	// Only apply defaults if CORS is enabled (either AllowAllOrigins or specific origins)
	isEnabled := (SsmConfig.Configuration.CORS.AllowAllOrigins != nil && *SsmConfig.Configuration.CORS.AllowAllOrigins) ||
		(len(SsmConfig.Configuration.CORS.AllowOrigins) > 0)

	if !isEnabled {
		logger.CfgLog.Info("CORS is disabled. No origins configured")
		return
	}

	logger.CfgLog.Info("CORS is enabled. Applying secure defaults for unset fields")

	// Set default allowed methods (standard safe HTTP methods)
	if len(SsmConfig.Configuration.CORS.AllowMethods) == 0 {
		SsmConfig.Configuration.CORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
		logger.CfgLog.Infof("CORS AllowMethods not set. Using default: %v", SsmConfig.Configuration.CORS.AllowMethods)
	}

	// Set default allowed headers (common safe headers)
	if len(SsmConfig.Configuration.CORS.AllowHeaders) == 0 {
		SsmConfig.Configuration.CORS.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
		logger.CfgLog.Infof("CORS AllowHeaders not set. Using default: %v", SsmConfig.Configuration.CORS.AllowHeaders)
	}

	// Set default exposed headers (safe to expose)
	if len(SsmConfig.Configuration.CORS.ExposeHeaders) == 0 {
		SsmConfig.Configuration.CORS.ExposeHeaders = []string{"Content-Length"}
		logger.CfgLog.Infof("CORS ExposeHeaders not set. Using default: %v", SsmConfig.Configuration.CORS.ExposeHeaders)
	}

	// Set AllowCredentials to false by default for security (only enable with specific origins)
	if SsmConfig.Configuration.CORS.AllowCredentials == nil {
		credentialsDisabled := false
		// Only allow credentials if specific origins are set (not AllowAllOrigins)
		if SsmConfig.Configuration.CORS.AllowAllOrigins != nil && !*SsmConfig.Configuration.CORS.AllowAllOrigins {
			credentialsDisabled = true // Can be enabled with specific origins
		}
		SsmConfig.Configuration.CORS.AllowCredentials = &credentialsDisabled
		logger.CfgLog.Infof("CORS AllowCredentials not set. Defaulting to %v", credentialsDisabled)
	}

	// Validate: If AllowAllOrigins is true, credentials must be false
	if SsmConfig.Configuration.CORS.AllowAllOrigins != nil && *SsmConfig.Configuration.CORS.AllowAllOrigins {
		if SsmConfig.Configuration.CORS.AllowCredentials != nil && *SsmConfig.Configuration.CORS.AllowCredentials {
			credentialsDisabled := false
			SsmConfig.Configuration.CORS.AllowCredentials = &credentialsDisabled
			logger.CfgLog.Warn("CORS AllowCredentials cannot be true when AllowAllOrigins is true. Setting to false")
		}
	}

	// Set MaxAge to 12 hours by default (standard recommendation)
	if SsmConfig.Configuration.CORS.MaxAge == nil {
		defaultMaxAge := int(12 * time.Hour / time.Second) // 43200 seconds
		SsmConfig.Configuration.CORS.MaxAge = &defaultMaxAge
		logger.CfgLog.Infof("CORS MaxAge not set. Using default: %d seconds (12 hours)", defaultMaxAge)
	}

	// Set AllowWildcard to false by default for security
	if SsmConfig.Configuration.CORS.AllowWildcard == nil {
		wildcardDisabled := false
		SsmConfig.Configuration.CORS.AllowWildcard = &wildcardDisabled
		logger.CfgLog.Info("CORS AllowWildcard not set. Defaulting to false for security")
	}

	// Set AllowBrowserExtensions to false by default
	if SsmConfig.Configuration.CORS.AllowBrowserExtensions == nil {
		extensionsDisabled := false
		SsmConfig.Configuration.CORS.AllowBrowserExtensions = &extensionsDisabled
		logger.CfgLog.Info("CORS AllowBrowserExtensions not set. Defaulting to false")
	}

	// Set AllowWebSockets to false by default
	if SsmConfig.Configuration.CORS.AllowWebSockets == nil {
		websocketsDisabled := false
		SsmConfig.Configuration.CORS.AllowWebSockets = &websocketsDisabled
		logger.CfgLog.Info("CORS AllowWebSockets not set. Defaulting to false")
	}

	// Set AllowFiles to false by default (security risk)
	if SsmConfig.Configuration.CORS.AllowFiles == nil {
		filesDisabled := false
		SsmConfig.Configuration.CORS.AllowFiles = &filesDisabled
		logger.CfgLog.Info("CORS AllowFiles not set. Defaulting to false for security")
	}

	// Set AllowPrivateNetwork to false by default
	if SsmConfig.Configuration.CORS.AllowPrivateNetwork == nil {
		privateNetworkDisabled := false
		SsmConfig.Configuration.CORS.AllowPrivateNetwork = &privateNetworkDisabled
		logger.CfgLog.Info("CORS AllowPrivateNetwork not set. Defaulting to false")
	}

	// Set default OPTIONS response status code to 204 (No Content)
	if SsmConfig.Configuration.CORS.OptionsResponseStatusCode == nil {
		defaultOptionsStatus := 204
		SsmConfig.Configuration.CORS.OptionsResponseStatusCode = &defaultOptionsStatus
		logger.CfgLog.Info("CORS OptionsResponseStatusCode not set. Defaulting to 204")
	}

	// Initialize CustomSchemas as empty if nil
	if SsmConfig.Configuration.CORS.CustomSchemas == nil {
		SsmConfig.Configuration.CORS.CustomSchemas = []string{}
		logger.CfgLog.Info("CORS CustomSchemas not set. Using empty list")
	}

	logger.CfgLog.Info("CORS configuration initialized with secure defaults")
}
