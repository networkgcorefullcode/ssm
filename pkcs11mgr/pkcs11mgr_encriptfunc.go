package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

// EncryptWithAESKey performs encryption using a key object already in the token.
// NOTE: parámetros específicos del mecanismo (p.ej. GCM params) pueden necesitar ajustar según tu módulo.
func (m *Manager) EncryptKey(keyHandle pkcs11.ObjectHandle, iv, plaintext []byte, encryptALGORITHM uint) ([]byte, error) {
	logger.AppLog.Infof("Encrypting data with key handle=%v, algorithm=%v (0x%X)", keyHandle, encryptALGORITHM, encryptALGORITHM)
	logger.AppLog.Infof("IV length: %d bytes, Plaintext length: %d bytes", len(iv), len(plaintext))

	// Create mechanism with IV parameter
	mech := pkcs11.NewMechanism(encryptALGORITHM, iv)
	logger.AppLog.Infof("Mechanism created successfully")

	if err := m.ctx.EncryptInit(m.session, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("EncryptInit failed for mechanism 0x%X: %v", encryptALGORITHM, err)
		return nil, err
	}
	out, err := m.ctx.Encrypt(m.session, plaintext)
	if err != nil {
		logger.AppLog.Errorf("Encrypt failed: %v", err)
		return nil, err
	}
	logger.AppLog.Infof("Encryption successful, ciphertext length=%d", len(out))
	return out, nil
}

func (m *Manager) DecryptKey(keyHandle pkcs11.ObjectHandle, iv, ciphertext []byte, decriptALGORITHM uint) ([]byte, error) {
	logger.AppLog.Infof("Decrypting data with key handle=%v, algorithm=%v", keyHandle, decriptALGORITHM)
	mech := pkcs11.NewMechanism(decriptALGORITHM, iv)
	if err := m.ctx.DecryptInit(m.session, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("DecryptInit failed: %v", err)
		return nil, err
	}
	out, err := m.ctx.Decrypt(m.session, ciphertext)
	if err != nil {
		logger.AppLog.Errorf("Decrypt failed: %v", err)
		return nil, err
	}
	logger.AppLog.Infof("Decryption successful, plaintext length=%d", len(out))
	return out, nil
}
