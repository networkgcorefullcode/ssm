package constants

// Package constants holds application-wide constant values and mappings.
const (
	APP_NAME         = "ssm"
	APP_VERSION      = "0.3.0"
	APP_AUTHOR       = "NetworkG"
	APP_EMAIL        = "support@networkg.com"
	APP_URL          = "https://networkg.com"
	APP_DESC         = "Simple Secret Manager - SSM"
	APP_PATH_SWAGGER = "handlers/api/"
)

var AlgorithmLabelMap = map[int]string{
	ALGORITM_AES256:          LABEL_K4_KEY_AES,
	ALGORITM_AES128:          LABEL_K4_KEY_AES,
	ALGORITM_DES:             LABEL_K4_KEY_DES,
	ALGORITM_DES3:            LABEL_K4_KEY_DES3,
	ALGORITM_AES256_OurUsers: LABEL_ENCRIPTION_KEY_AES256,
	ALGORITM_AES128_OurUsers: LABEL_ENCRIPTION_KEY_AES128,
	ALGORITM_DES_OurUsers:    LABEL_ENCRIPTION_KEY_DES,
	ALGORITM_DES3_OurUsers:   LABEL_ENCRIPTION_KEY_DES3,
}

const (
	LABEL_ENCRYPTION_KEY        = "SSM_ENC_KEY"
	LABEL_K4_KEY_AES            = "K4_AES"
	LABEL_K4_KEY_DES            = "K4_DES"
	LABEL_K4_KEY_DES3           = "K4_DES3"
	LABEL_ENCRIPTION_KEY_AES256 = "KEY_ENCRIPTION_AES256"
	LABEL_ENCRIPTION_KEY_AES128 = "KEY_ENCRIPTION_AES128"
	LABEL_ENCRIPTION_KEY_DES    = "KEY_ENCRIPTION_DES"
	LABEL_ENCRIPTION_KEY_DES3   = "KEY_ENCRIPTION_DES3"
	ALGORITM_AES256             = 1
	ALGORITM_AES128             = 2
	ALGORITM_DES                = 3
	ALGORITM_DES3               = 4
	ALGORITM_AES256_OurUsers    = 5
	ALGORITM_AES128_OurUsers    = 6
	ALGORITM_DES_OurUsers       = 7
	ALGORITM_DES3_OurUsers      = 8
)

var KeyLabelsExternalAllow [3]string = [3]string{
	LABEL_K4_KEY_AES,
	LABEL_K4_KEY_DES,
	LABEL_K4_KEY_DES3,
}

var KeyLabelsInternalAllow [4]string = [4]string{
	LABEL_ENCRIPTION_KEY_AES256,
	LABEL_ENCRIPTION_KEY_AES128,
	LABEL_ENCRIPTION_KEY_DES,
	LABEL_ENCRIPTION_KEY_DES3,
}

const (
	TYPE_AES  = "AES"
	TYPE_DES  = "DES"
	TYPE_DES3 = "DES3"
)

var KeyTypeAllow [3]string = [3]string{
	TYPE_AES,
	TYPE_DES,
	TYPE_DES3,
}
