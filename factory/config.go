/*
 * SSM Configuration Factory
 */

package factory

import (
	"github.com/omec-project/util/logger"
)

const (
	SSM_EXPECTED_CONFIG_VERSION = "1.0.0"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *logger.Logger `yaml:"logger"`
	CfgLocation   string
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

const (
	SSM_DEFAULT_IPV4     = "127.0.0.18"
	SSM_DEFAULT_PORT     = "8000"
	SSM_DEFAULT_PORT_INT = 8000
	SSM_DEFAULT_NRFURI   = "https://127.0.0.10:8000"
)

type Configuration struct {
	SsmName         string     `yaml:"ssmName,omitempty"`
	SsmId           string     `yaml:"ssmId,omitempty"`
	SocketPath      string     `yaml:"socketPath,omitempty"`
	PkcsPath        string     `yaml:"pkcsPath,omitempty"`
	Pin             string     `yaml:"pin,omitempty"`
	LotsNumber      int        `yaml:"lotsNumber,omitempty"`
	BindAddr        string     `yaml:"bindAddr,omitempty"`
	ExposeSwaggerUi *bool      `yaml:"exposeSwaggerUi,omitempty"`
	IsHttps         *bool      `yaml:"isHttps,omitempty"`
	CertFile        string     `yaml:"certFile,omitempty"`
	KeyFile         string     `yaml:"keyFile,omitempty"`
	MaxSessions     int        `yaml:"maxSessions,omitempty"`
	IsSecure        bool       `yaml:"isSecure,omitempty"`
	CORS            *CORS      `yaml:"cors,omitempty"`
	Mongodb         *Mongodb   `yaml:"mongodb"`
	RateLimit       *RateLimit `yaml:"rateLimit,omitempty"`
}

type Mongodb struct {
	Name   string `yaml:"name,omitempty"`
	Url    string `yaml:"url,omitempty"`
	DBName string `yaml:"dbName,omitempty"`
}

type RateLimit struct {
	Enabled         bool `yaml:"enabled,omitempty"`
	RequestsPerMin  int  `yaml:"requestsPerMin,omitempty"`
	BurstSize       int  `yaml:"burstSize,omitempty"`
	CleanupInterval int  `yaml:"cleanupInterval,omitempty"` // en minutos
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

// GetRateLimit returns rate limit configuration with defaults
func (c *Config) GetRateLimit() *RateLimit {
	if c.Configuration != nil && c.Configuration.RateLimit != nil {
		rl := c.Configuration.RateLimit

		// Set defaults if values are not configured
		if rl.RequestsPerMin <= 0 {
			rl.RequestsPerMin = 100 // default: 100 requests per minute
		}
		if rl.BurstSize <= 0 {
			rl.BurstSize = 10 // default: allow burst of 10 requests
		}
		if rl.CleanupInterval <= 0 {
			rl.CleanupInterval = 15 // default: cleanup every 15 minutes
		}

		return rl
	}

	// Return default configuration if none provided
	return &RateLimit{
		Enabled:         false, // disabled by default
		RequestsPerMin:  100,
		BurstSize:       10,
		CleanupInterval: 15,
	}
}
