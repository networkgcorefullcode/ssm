/*
 * SSM Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"regexp"

	"github.com/networkgcorefullcode/ssm/logger"
	"gopkg.in/yaml.v2"
)

var SsmConfig Config

const SsmFID_PATTERN = "^[A-Fa-f0-9]{6}$"

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	content, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(content, &SsmConfig); err != nil {
		return err
	}
	if SsmConfig.Configuration.SsmId == "" {
		SsmConfig.Configuration.SsmId = "cafe00"
		logger.CfgLog.Infof("ssmId not set in configuration file. Using %s", SsmConfig.Configuration.SsmId)
	}
	err = validateSsmId(SsmConfig.Configuration.SsmId)

	if SsmConfig.Configuration.BindAddr == "" {
		SSM_DEFAULT_BIND_ADDR := fmt.Sprintf("%s:%d", SSM_DEFAULT_IPV4, SSM_DEFAULT_PORT_INT)
		SsmConfig.Configuration.BindAddr = SSM_DEFAULT_BIND_ADDR
		logger.CfgLog.Infof("bindAddr not set in configuration file. Using %s", SsmConfig.Configuration.BindAddr)
	}
	if SsmConfig.Configuration.PkcsPath == "" {
		SsmConfig.Configuration.PkcsPath = "/usr/local/lib/softhsm/libsofthsm2.so"
		logger.CfgLog.Infof("pkcsPath not set in configuration file. Using %s", SsmConfig.Configuration.PkcsPath)
	}

	if SsmConfig.Configuration.Pin == "" {
		SsmConfig.Configuration.Pin = "1234"
		logger.CfgLog.Infof("pin not set in configuration file. Using default pin")
	}

	if SsmConfig.Configuration.MaxSessions == 0 {
		SsmConfig.Configuration.MaxSessions = 10
		logger.CfgLog.Infof("maxSessions not set in configuration file. Using default value: %d", SsmConfig.Configuration.MaxSessions)
	}

	// Check CORS configs and set default if not present
	initializeCORSConfig()

	return err
}

func CheckConfigVersion() error {
	currentVersion := SsmConfig.GetVersion()

	if currentVersion != SSM_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s]",
			currentVersion, SSM_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}

func validateSsmId(ssmId string) error {
	ssmIdMatch, err := regexp.MatchString(SsmFID_PATTERN, ssmId)
	if err != nil {
		return fmt.Errorf("invalid ssmId: %s. It should match the following pattern: `%s`", ssmId, SsmFID_PATTERN)
	}
	if !ssmIdMatch {
		return fmt.Errorf("invalid ssmId: %s. It should match the following pattern: `%s`", ssmId, SsmFID_PATTERN)
	}
	return nil
}
