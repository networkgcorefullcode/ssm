package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

// GetMechanismInfo retrieves information about a specific mechanism
func (m *Manager) GetMechanismInfo(mechanismType uint) (*pkcs11.MechanismInfo, error) {
	info, err := m.ctx.GetMechanismInfo(m.slot, []*pkcs11.Mechanism{pkcs11.NewMechanism(mechanismType, nil)})
	if err != nil {
		logger.AppLog.Errorf("Failed to get mechanism info for 0x%X: %v", mechanismType, err)
		return nil, err
	}
	return info, nil
}

// ListSupportedMechanisms lists all mechanisms supported by the slot
func (m *Manager) ListSupportedMechanisms() ([]uint, error) {
	mechanisms, err := m.ctx.GetMechanismList(m.slot)
	if err != nil {
		logger.AppLog.Errorf("Failed to get mechanism list: %v", err)
		return nil, err
	}

	logger.AppLog.Infof("Slot supports %d mechanisms", len(mechanisms))

	// Log some important ones
	importantMechs := map[uint]string{
		pkcs11.CKM_AES_KEY_GEN:  "CKM_AES_KEY_GEN",
		pkcs11.CKM_AES_CBC:      "CKM_AES_CBC",
		pkcs11.CKM_AES_CBC_PAD:  "CKM_AES_CBC_PAD",
		pkcs11.CKM_DES_KEY_GEN:  "CKM_DES_KEY_GEN",
		pkcs11.CKM_DES_CBC:      "CKM_DES_CBC",
		pkcs11.CKM_DES_CBC_PAD:  "CKM_DES_CBC_PAD",
		pkcs11.CKM_DES_ECB:      "CKM_DES_ECB",
		pkcs11.CKM_DES3_KEY_GEN: "CKM_DES3_KEY_GEN",
		pkcs11.CKM_DES3_CBC:     "CKM_DES3_CBC",
		pkcs11.CKM_DES3_CBC_PAD: "CKM_DES3_CBC_PAD",
		pkcs11.CKM_DES3_ECB:     "CKM_DES3_ECB",
	}

	for _, mech := range mechanisms {
		if name, ok := importantMechs[mech]; ok {
			logger.AppLog.Infof("Mechanism supported: %s (0x%X)", name, mech)
		}
	}

	return mechanisms, nil
}

// IsMechanismSupported checks if a mechanism is supported
func (m *Manager) IsMechanismSupported(mechanismType uint) bool {
	mechanisms, err := m.ctx.GetMechanismList(m.slot)
	if err != nil {
		return false
	}

	for _, mech := range mechanisms {
		if mech == mechanismType {
			return true
		}
	}
	return false
}
