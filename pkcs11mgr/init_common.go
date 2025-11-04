package pkcs11mgr

import (
	"time"

	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
)

var mgr *Manager

func SetPKCS11Manager(manager *Manager) {
	mgr = manager
}

func InitPKCS11() {
	session := mgr.GetSession()
	defer mgr.LogoutSession(session)

	if err := InitAuditKey(session); err != nil {
		logger.AppLog.Errorf("Failed to initialize audit key: %v", err)
		return
	}

	time.Sleep(time.Second * 2)

	if err := InitJWTKey(session); err != nil {
		logger.AppLog.Errorf("Failed to initialize jwt key: %v", err)
		return
	}

	_, err := FindKey(constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256, 0, *session)
	if err != nil && err.Error() == constants.ERROR_STRING_KEY_NOT_FOUND {
		_, _, err = GenerateAESKey(constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256, 0, 256, *session)
		if err != nil {
			logger.AppLog.Errorf("Failed to generate AES key: %v", err)
			return
		}
	}

	if err != nil {
		logger.AppLog.Errorf("Error during PKCS11 initialization: %v", err)
	}
}
