package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

// DeleteKey deletes a key from the HSM by its handle
func (m *Manager) DeleteKey(handle pkcs11.ObjectHandle) error {
	logger.AppLog.Infof("Deleting key with handle: %v", handle)

	if err := m.ctx.DestroyObject(m.session, handle); err != nil {
		logger.AppLog.Errorf("Failed to delete key: %v", err)
		return err
	}

	logger.AppLog.Infof("Key deleted successfully")
	return nil
}

// DeleteKeyByLabel deletes all keys with the specified label
func (m *Manager) DeleteKeyByLabel(label string) error {
	logger.AppLog.Infof("Deleting keys with label: %s", label)

	handles, err := m.FindKeysLabel(label)
	if err != nil {
		logger.AppLog.Errorf("Failed to find keys: %v", err)
		return err
	}

	if len(handles) == 0 {
		logger.AppLog.Warnf("No keys found with label: %s", label)
		return nil
	}

	for _, handle := range handles {
		if err := m.DeleteKey(handle); err != nil {
			logger.AppLog.Errorf("Failed to delete key handle %v: %v", handle, err)
			return err
		}
	}

	logger.AppLog.Infof("Deleted %d keys with label: %s", len(handles), label)
	return nil
}
