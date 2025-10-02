/*
 * SSM Configuration Factory
 */

package factory

import (
	"fmt"
	"net/url"
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
		logger.CfgLog.Infof("amfId not set in configuration file. Using %s", SsmConfig.Configuration.SsmId)
	}
	if SsmConfig.Configuration.WebuiUri == "" {
		SsmConfig.Configuration.WebuiUri = "http://webui:5001"
		logger.CfgLog.Infof("webuiUri not set in configuration file. Using %s", SsmConfig.Configuration.WebuiUri)
	}
	if SsmConfig.Configuration.KafkaInfo.EnableKafka == nil {
		enableKafka := true
		SsmConfig.Configuration.KafkaInfo.EnableKafka = &enableKafka
	}
	// if SsmConfig.Configuration.Telemetry != nil && SsmConfig.Configuration.Telemetry.Enabled {
	// 	if SsmConfig.Configuration.Telemetry.Ratio == nil {
	// 		defaultRatio := 1.0
	// 		SsmConfig.Configuration.Telemetry.Ratio = &defaultRatio
	// 	}

	// 	if SsmConfig.Configuration.Telemetry.OtlpEndpoint == "" {
	// 		return fmt.Errorf("OTLP endpoint is not set in the configuration")
	// 	}
	// }
	if err = validateWebuiUri(SsmConfig.Configuration.WebuiUri); err != nil {
		return err
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

func validateWebuiUri(uri string) error {
	parsedUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return fmt.Errorf("unsupported scheme for webuiUri: %s", parsedUrl.Scheme)
	}
	if parsedUrl.Hostname() == "" {
		return fmt.Errorf("missing host in webuiUri")
	}
	return nil
}

func validateSsmId(amfId string) error {
	amfIdMatch, err := regexp.MatchString(SsmFID_PATTERN, amfId)
	if err != nil {
		return fmt.Errorf("invalid amfId: %s. It should match the following pattern: `%s`", amfId, SsmFID_PATTERN)
	}
	if !amfIdMatch {
		return fmt.Errorf("invalid amfId: %s. It should match the following pattern: `%s`", amfId, SsmFID_PATTERN)
	}
	return nil
}
