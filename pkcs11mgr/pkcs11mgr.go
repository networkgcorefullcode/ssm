package pkcs11mgr

import (
	"errors"

	"github.com/miekg/pkcs11"
	ssm_consts "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
)

type Manager struct {
	ctx     *pkcs11.Ctx
	slot    uint
	session pkcs11.SessionHandle
	pin     string
}

func New(modulePath string, slot uint, pin string) (*Manager, error) {
	ctx := pkcs11.New(modulePath)
	if ctx == nil {
		return nil, errors.New("pkcs11.New returned nil")
	}
	if err := ctx.Initialize(); err != nil {
		return nil, err
	}
	mgr := &Manager{ctx: ctx, slot: slot, pin: pin}
	logger.AppLog.Infoln("PKCS#11 module initialized")
	return mgr, nil
}

// Open a Session to operate with the SSM
func (m *Manager) OpenSession() error {
	// Open a session with the specified slot
	logger.AppLog.Infoln("Opening PKCS#11 session")
	session, err := m.ctx.OpenSession(m.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return err
	}
	m.session = session
	if err := m.ctx.Login(m.session, pkcs11.CKU_USER, m.pin); err != nil {
		return err
	}
	logger.AppLog.Infoln("PKCS#11 session has been opened")
	return nil
}

// CloseSession logs out and closes the session
func (m *Manager) CloseSession() {
	logger.AppLog.Infoln("Closing PKCS#11 session")
	if m.session != 0 {
		_ = m.ctx.Logout(m.session)
		_ = m.ctx.CloseSession(m.session)
		m.session = 0
	}
	logger.AppLog.Infoln("PKCS#11 session has been closed")
}

// Finalize cleans up the PKCS#11 context
func (m *Manager) Finalize() {
	m.CloseSession()
	if m.ctx != nil {
		_ = m.ctx.Finalize()
		m.ctx.Destroy()
		m.ctx = nil
	}
}

// GenerateAESKey creates an AES key object inside SoftHSM and returns its object handle (as uint)
func (m *Manager) GenerateAESKey(label string, id []byte, bits int) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Generating AES key: label=%s, bits=%d", label, bits)
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, bits/8),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true), // store persistently in token
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}
	handle, err := m.ctx.GenerateKey(m.session, []*pkcs11.Mechanism{mech}, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate AES key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("AES key generated successfully: handle=%v", handle)
	return handle, nil
}

// EncryptWithAESKey performs encryption using a key object already in the token.
// NOTE: parámetros específicos del mecanismo (p.ej. GCM params) pueden necesitar ajustar según tu módulo.
func (m *Manager) EncryptKey(keyHandle pkcs11.ObjectHandle, iv, plaintext []byte, encryptAlgoritm uint) ([]byte, error) {
	logger.AppLog.Infof("Encrypting data with key handle=%v, algorithm=%v", keyHandle, encryptAlgoritm)
	mech := pkcs11.NewMechanism(encryptAlgoritm, iv)
	if err := m.ctx.EncryptInit(m.session, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		logger.AppLog.Errorf("EncryptInit failed: %v", err)
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

func (m *Manager) DecryptKey(keyHandle pkcs11.ObjectHandle, iv, ciphertext []byte, decriptAlgoritm uint) ([]byte, error) {
	logger.AppLog.Infof("Decrypting data with key handle=%v, algorithm=%v", keyHandle, decriptAlgoritm)
	mech := pkcs11.NewMechanism(decriptAlgoritm, iv)
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

// StoreKey creates a key object inside SoftHSM from raw key bytes and returns its object handle
func (m *Manager) StoreKey(label string, key []byte, id []byte, keyType string) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Storing key: label=%s, keyType=%s, keyLen=%d", label, keyType, len(key))
	var keyTypeuint uint
	switch keyType {
	case ssm_consts.TYPE_AES:
		keyTypeuint = pkcs11.CKK_AES
	case ssm_consts.TYPE_DES3:
		keyTypeuint = pkcs11.CKK_DES3
	case ssm_consts.TYPE_DES:
		keyTypeuint = pkcs11.CKK_DES
	default:
		return 0, errors.New("unsupported key type")
	}
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, keyTypeuint),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, key),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, false),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}

	// Check if key already exists before creating it
	existingHandle, err := m.FindKey(label, string(id))
	if err == nil && existingHandle != 0 {
		logger.AppLog.Infof("Key with label '%s' already exists, returning existing handle: %v", label, existingHandle)
		return existingHandle, nil
	}

	handle, err := m.ctx.CreateObject(m.session, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to store key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("Key stored successfully: handle=%v", handle)
	return handle, nil
}

// FindKey returns the object handle for a given label, or 0 if not found
func (m *Manager) FindKey(label, id string) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if id != "" {
		template = append(template, pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(id)))
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

// DeleteKey removes a key object from the HSM by label and optionally by ID
func (m *Manager) DeleteKey(label, id string) error {
	logger.AppLog.Infof("Attempting to delete key with label: %s", label)

	// Find the key first
	handle, err := m.FindKey(label, id)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key for deletion: %v", err)
		return err
	}

	// If key not found (handle is 0), return error
	if handle == 0 {
		logger.AppLog.Errorf("Key with label '%s' not found for deletion", label)
		return errors.New("key not found")
	}

	// Delete the key object
	if err := m.ctx.DestroyObject(m.session, handle); err != nil {
		logger.AppLog.Errorf("Failed to delete key with handle %v: %v", handle, err)
		return err
	}

	logger.AppLog.Infof("Key successfully deleted: label=%s, handle=%v", label, handle)
	return nil
}

// DeleteAllKeys removes all secret key objects from the current session/slot
// WARNING: This will delete ALL keys in the token - use with caution!
func (m *Manager) DeleteAllKeys() error {
	logger.AppLog.Warnln("Attempting to delete ALL keys from the HSM token")

	// Search for all secret key objects
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed for delete all: %v", err)
		return err
	}
	defer m.ctx.FindObjectsFinal(m.session)

	// Find all matching objects (up to 1000 keys)
	handles, _, err := m.ctx.FindObjects(m.session, 1000)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed for delete all: %v", err)
		return err
	}

	if len(handles) == 0 {
		logger.AppLog.Infoln("No keys found to delete")
		return nil
	}

	logger.AppLog.Infof("Found %d keys to delete", len(handles))

	// Delete each key object
	deletedCount := 0
	for _, handle := range handles {
		if err := m.ctx.DestroyObject(m.session, handle); err != nil {
			logger.AppLog.Errorf("Failed to delete key with handle %v: %v", handle, err)
			// Continue with other keys even if one fails
			continue
		}
		deletedCount++
		logger.AppLog.Infof("Deleted key with handle: %v", handle)
	}

	logger.AppLog.Infof("Successfully deleted %d out of %d keys", deletedCount, len(handles))

	if deletedCount != len(handles) {
		return errors.New("not all keys were successfully deleted")
	}

	return nil
}

// UpdateKey updates an existing key by deleting the old one and creating a new one with the updated value
// This function combines delete and store operations to effectively "update" a key
func (m *Manager) UpdateKey(label string, newKeyValue []byte, id []byte, keyType string) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Updating key: label=%s, keyType=%s, newKeyLen=%d", label, keyType, len(newKeyValue))

	// First, check if the key exists
	existingHandle, err := m.FindKey(label, string(id))
	if err != nil {
		logger.AppLog.Errorf("Error searching for key to update: %v", err)
		return 0, err
	}

	if existingHandle == 0 {
		logger.AppLog.Errorf("Key with label '%s' not found for update", label)
		return 0, errors.New("key not found for update")
	}

	logger.AppLog.Infof("Found existing key to update: handle=%v", existingHandle)

	// Delete the existing key
	if err := m.DeleteKey(label, string(id)); err != nil {
		logger.AppLog.Errorf("Failed to delete existing key for update: %v", err)
		return 0, err
	}

	logger.AppLog.Infof("Existing key deleted, creating new key with updated value")

	newHandle, err := m.StoreKey(label, newKeyValue, id, keyType)
	if err != nil {
		logger.AppLog.Errorf("Failed to create updated key: %v", err)
		return 0, err
	}

	logger.AppLog.Infof("Key updated successfully: label=%s, newHandle=%v", label, newHandle)
	return newHandle, nil
}
