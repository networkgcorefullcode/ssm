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
