package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/utils"
)

// FindKey returns the object handle for a given label, or 0 if not found
func (m *Manager) FindKey(label string, id int32) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if id != 0 {
		// Convert int32 to byte slice
		b := utils.Int32ToByte(id)
		template = append(template, pkcs11.NewAttribute(pkcs11.CKA_ID, b))
	}

	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return 0, err
	}
	defer m.ctx.FindObjectsFinal(m.session)

	handles, _, err := m.ctx.FindObjects(m.session, 1)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed: %v", err)
		return 0, err
	}
	if len(handles) == 0 {
		logger.AppLog.Warnf("No key found with label: %s", label)
		return 0, err
	}
	logger.AppLog.Infof("Key found: handle=%v", handles[0])
	return handles[0], nil
}
