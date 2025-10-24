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
	SsmName         string `yaml:"ssmName,omitempty"`
	SsmId           string `yaml:"ssmId,omitempty"`
	SocketPath      string `yaml:"socketPath,omitempty"`
	PkcsPath        string `yaml:"pkcsPath,omitempty"`
	Pin             string `yaml:"pin,omitempty"`
	LotsNumber      int    `yaml:"lotsNumber,omitempty"`
	BindAddr        string `yaml:"bindAddr,omitempty"`
	ExposeSwaggerUi *bool  `yaml:"exposeSwaggerUi,omitempty"`
	IsHttps         *bool  `yaml:"isHttps,omitempty"`
	CertFile        string `yaml:"certFile,omitempty"`
	KeyFile         string `yaml:"keyFile,omitempty"`
	// HandlersPoolConect bool   `yaml:"handlersPoolConect,omitempty"`
	// PoolConfig         *pkcs11mgr.PoolConfig `yaml:"poolConfig,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
