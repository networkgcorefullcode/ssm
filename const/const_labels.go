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
	ALGORITHM_AES256:          LABEL_K4_KEY_AES,
	ALGORITHM_AES128:          LABEL_K4_KEY_AES,
	ALGORITHM_DES:             LABEL_K4_KEY_DES,
	ALGORITHM_DES3:            LABEL_K4_KEY_DES3,
	ALGORITHM_AES256_OurUsers: LABEL_ENCRYPTION_KEY_AES256,
	ALGORITHM_AES128_OurUsers: LABEL_ENCRYPTION_KEY_AES128,
	ALGORITHM_DES_OurUsers:    LABEL_ENCRYPTION_KEY_DES,
	ALGORITHM_DES3_OurUsers:   LABEL_ENCRYPTION_KEY_DES3,
}

var LabelAlgorithmMap = map[string]int{
	LABEL_K4_KEY_AES:            ALGORITHM_AES256,
	LABEL_K4_KEY_DES:            ALGORITHM_DES,
	LABEL_K4_KEY_DES3:           ALGORITHM_DES3,
	LABEL_ENCRYPTION_KEY_AES256: ALGORITHM_AES256_OurUsers,
	LABEL_ENCRYPTION_KEY_AES128: ALGORITHM_AES128_OurUsers,
	LABEL_ENCRYPTION_KEY_DES:    ALGORITHM_DES_OurUsers,
	LABEL_ENCRYPTION_KEY_DES3:   ALGORITHM_DES3_OurUsers,
}

const (
	LABEL_ENCRYPTION_KEY                 = "SSM_ENC_KEY"
	LABEL_K4_KEY_AES                     = "K4_AES"
	LABEL_K4_KEY_DES                     = "K4_DES"
	LABEL_K4_KEY_DES3                    = "K4_DES3"
	LABEL_ENCRYPTION_KEY_AES256          = "KEY_ENCRYPTION_AES256"
	LABEL_ENCRYPTION_KEY_AES128          = "KEY_ENCRYPTION_AES128"
	LABEL_ENCRYPTION_KEY_DES             = "KEY_ENCRYPTION_DES"
	LABEL_ENCRYPTION_KEY_DES3            = "KEY_ENCRYPTION_DES3"
	AuditKeyLabel                        = "AUDIT_SIGNING_KEY"
	LABEL_ENCRYPTION_KEY_INTERNAL_AES256 = "ENCRYPTION_KEY_INTERNAL_AES256"
	LABEL_ENCRYPTION_KEY_INTERNAL_AES128 = "ENCRYPTION_KEY_INTERNAL_AES128"
	ALGORITHM_AES256                     = 1
	ALGORITHM_AES128                     = 2
	ALGORITHM_DES                        = 3
	ALGORITHM_DES3                       = 4
	ALGORITHM_AES256_OurUsers            = 5
	ALGORITHM_AES128_OurUsers            = 6
	ALGORITHM_DES_OurUsers               = 7
	ALGORITHM_DES3_OurUsers              = 8
)

var KeyLabelsExternalAllow [3]string = [3]string{
	LABEL_K4_KEY_AES,
	LABEL_K4_KEY_DES,
	LABEL_K4_KEY_DES3,
}

var KeyLabelsInternalAllow [4]string = [4]string{
	LABEL_ENCRYPTION_KEY_AES256,
	LABEL_ENCRYPTION_KEY_AES128,
	LABEL_ENCRYPTION_KEY_DES,
	LABEL_ENCRYPTION_KEY_DES3,
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
