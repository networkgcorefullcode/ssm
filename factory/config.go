/*
 * SSM Configuration Factory
 */

package factory

import (
	"time"

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
	Rcvd          bool
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

type Mongodb struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type KafkaInfo struct {
	EnableKafka *bool  `yaml:"enableKafka,omitempty"`
	BrokerUri   string `yaml:"brokerUri,omitempty"`
	BrokerPort  int    `yaml:"brokerPort,omitempty"`
	Topic       string `yaml:"topicName,omitempty"`
}

// type TelemetryConfig struct {
// 	Enabled      bool     `yaml:"enabled,omitempty"`       // Optional; defaults to false
// 	OtlpEndpoint string   `yaml:"otlp_endpoint,omitempty"` // Mandatory if enabled=true
// 	Ratio        *float64 `yaml:"ratio,omitempty"`         // Optional; defaults to 1.0
// }

type Configuration struct {
	SsmName         string      `yaml:"ssmName,omitempty"`
	SsmId           string      `yaml:"ssmId,omitempty"`
	SocketPath      string      `yaml:"socketPath,omitempty"`
	SsmDBName       string      `yaml:"ssmDBName,omitempty"`
	Mongodb         *Mongodb    `yaml:"mongodb,omitempty"`
	Sbi             *Sbi        `yaml:"sbi,omitempty"`
	ServiceNameList []string    `yaml:"serviceNameList,omitempty"`
	SupportDnnList  []string    `yaml:"supportDnnList,omitempty"`
	NrfUri          string      `yaml:"nrfUri,omitempty"`
	WebuiUri        string      `yaml:"webuiUri"`
	Security        *Security   `yaml:"security,omitempty"`
	NetworkName     NetworkName `yaml:"networkName,omitempty"`

	EnableDbStore            bool      `yaml:"enableDBStore"`
	EnableNrfCaching         bool      `yaml:"enableNrfCaching"`
	NrfCacheEvictionInterval int       `yaml:"nrfCacheEvictionInterval,omitempty"`
	KafkaInfo                KafkaInfo `yaml:"kafkaInfo,omitempty"`
	DebugProfilePort         int       `yaml:"debugProfilePort,omitempty"`
}

type NetworkFeatureSupport5GS struct {
	Enable  bool  `yaml:"enable"`
	ImsVoPS uint8 `yaml:"imsVoPS"`
	Emc     uint8 `yaml:"emc"`
	Emf     uint8 `yaml:"emf"`
	IwkN26  uint8 `yaml:"iwkN26"`
	Mpsi    uint8 `yaml:"mpsi"`
	EmcN3   uint8 `yaml:"emcN3"`
	Mcsi    uint8 `yaml:"mcsi"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	TLS          *TLS   `yaml:"tls"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	BindingIPv4  string `yaml:"bindingIPv4,omitempty"`  // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty"`
}

type TLS struct {
	PEM string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

type Security struct {
	IntegrityOrder []string `yaml:"integrityOrder,omitempty"`
	CipheringOrder []string `yaml:"cipheringOrder,omitempty"`
}

type NetworkName struct {
	Full  string `yaml:"full"`
	Short string `yaml:"short,omitempty"`
}

type TimerValue struct {
	Enable        bool          `yaml:"enable"`
	ExpireTime    time.Duration `yaml:"expireTime"`
	MaxRetryTimes int           `yaml:"maxRetryTimes,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
