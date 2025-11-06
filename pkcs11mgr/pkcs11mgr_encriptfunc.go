package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

// EncryptWithAESKey performs encryption using a key object already in the token.
// NOTE: parámetros específicos del mecanismo (p.ej. GCM params) pueden necesitar ajustar según tu módulo.
func EncryptKey(keyHandle pkcs11.ObjectHandle, iv, plaintext []byte, encryptALGORITHM uint, s Session) ([]byte, error) {
	logger.AppLog.Infof("Encrypting data with key handle=%v, algorithm=%v (0x%X)", keyHandle, encryptALGORITHM, encryptALGORITHM)
	logger.AppLog.Infof("IV length: %d bytes, Plaintext length: %d bytes", len(iv), len(plaintext))

	// Create mechanism with IV parameter
	mech := pkcs11.NewMechanism(encryptALGORITHM, iv)
	logger.AppLog.Infof("Mechanism created successfully")

	if err := s.Ctx.EncryptInit(s.Handle, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("EncryptInit failed for mechanism 0x%X: %v", encryptALGORITHM, err)
		return nil, err
	}
	out, err := s.Ctx.Encrypt(s.Handle, plaintext)
	if err != nil {
		logger.AppLog.Errorf("Encrypt failed: %v", err)
		return nil, err
	}
	logger.AppLog.Infof("Encryption successful, ciphertext length=%d", len(out))
	return out, nil
}

// EncryptKeyAesGCM performs AES-GCM encryption using a key object already in the token.
// The function uses GCM parameters: IV (12 bytes recommended), AAD (optional), and Tag size (128 bits).
// Returns ciphertext with authentication tag appended at the end.
func EncryptKeyAesGCM(keyHandle pkcs11.ObjectHandle, iv, plaintext, aad []byte, s Session) ([]byte, error) {
	logger.AppLog.Infof("Encrypting data with AES-GCM, key handle=%v", keyHandle)
	logger.AppLog.Infof("IV length: %d bytes, Plaintext length: %d bytes, AAD length: %d bytes", len(iv), len(plaintext), len(aad))

	// Validate IV length (recommended 12 bytes for GCM)
	if len(iv) != 12 && len(iv) != 16 {
		logger.AppLog.Warnf("IV length is %d bytes. Recommended: 12 bytes for optimal GCM performance", len(iv))
	}

	// Create GCM parameters structure
	gcmParams := pkcs11.NewGCMParams(iv, aad, 128) // 128-bit tag (16 bytes)

	// Create mechanism with GCM parameters
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_GCM, gcmParams)
	logger.AppLog.Infof("GCM mechanism created successfully with tag size: 128 bits")

	if err := s.Ctx.EncryptInit(s.Handle, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("EncryptInit failed for AES-GCM: %v", err)
		return nil, err
	}

	// Encrypt returns ciphertext + authentication tag
	out, err := s.Ctx.Encrypt(s.Handle, plaintext)
	if err != nil {
		logger.AppLog.Errorf("AES-GCM Encrypt failed: %v", err)
		return nil, err
	}

	logger.AppLog.Infof("AES-GCM encryption successful, output length=%d (ciphertext + 16-byte tag)", len(out))
	return out, nil
}

func DecryptKey(keyHandle pkcs11.ObjectHandle, iv, ciphertext []byte, decriptALGORITHM uint, s Session) ([]byte, error) {
	logger.AppLog.Infof("Decrypting data with key handle=%v, algorithm=%v", keyHandle, decriptALGORITHM)
	mech := pkcs11.NewMechanism(decriptALGORITHM, iv)
	if err := s.Ctx.DecryptInit(s.Handle, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("DecryptInit failed: %v", err)
		return nil, err
	}
	out, err := s.Ctx.Decrypt(s.Handle, ciphertext)
	if err != nil {
		logger.AppLog.Errorf("Decrypt failed: %v", err)
		return nil, err
	}
	logger.AppLog.Infof("Decryption successful, plaintext length=%d", len(out))
	return out, nil
}

// DecryptKeyAesGCM performs AES-GCM decryption using a key object already in the token.
// The ciphertext should include the authentication tag appended at the end (last 16 bytes).
// The function uses GCM parameters: IV (12 bytes recommended), AAD (optional), and Tag size (128 bits).
// Returns plaintext if authentication succeeds, error otherwise.
func DecryptKeyAesGCM(keyHandle pkcs11.ObjectHandle, iv, ciphertext, aad []byte, s Session) ([]byte, error) {
	logger.AppLog.Infof("Decrypting data with AES-GCM, key handle=%v", keyHandle)
	logger.AppLog.Infof("IV length: %d bytes, Ciphertext+Tag length: %d bytes, AAD length: %d bytes", len(iv), len(ciphertext), len(aad))

	// Validate IV length
	if len(iv) != 12 && len(iv) != 16 {
		logger.AppLog.Warnf("IV length is %d bytes. Recommended: 12 bytes for optimal GCM performance", len(iv))
	}

	// Validate ciphertext length (must contain at least the tag)
	if len(ciphertext) < 16 {
		logger.AppLog.Errorf("Ciphertext too short: %d bytes (must include 16-byte authentication tag)", len(ciphertext))
		return nil, pkcs11.Error(pkcs11.CKR_ENCRYPTED_DATA_INVALID)
	}

	// Create GCM parameters structure
	gcmParams := pkcs11.NewGCMParams(iv, aad, 128) // 128-bit tag (16 bytes)

	// Create mechanism with GCM parameters
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_GCM, gcmParams)
	logger.AppLog.Infof("GCM mechanism created successfully with tag size: 128 bits")

	if err := s.Ctx.DecryptInit(s.Handle, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("DecryptInit failed for AES-GCM: %v", err)
		return nil, err
	}

	// Decrypt and verify authentication tag
	out, err := s.Ctx.Decrypt(s.Handle, ciphertext)
	if err != nil {
		logger.AppLog.Errorf("AES-GCM Decrypt failed (authentication may have failed): %v", err)
		return nil, err
	}

	logger.AppLog.Infof("AES-GCM decryption successful, plaintext length=%d", len(out))
	return out, nil
}
