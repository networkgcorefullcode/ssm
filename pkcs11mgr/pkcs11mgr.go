package pkcs11mgr

import (
	"errors"

	"github.com/miekg/pkcs11"
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
func (m *Manager) EncryptWithAESKey(keyHandle pkcs11.ObjectHandle, iv, plaintext []byte) ([]byte, error) {
	// Aquí debemos usar un mecanismo apropiado (p. ej. CKM_AES_GCM o CKM_AES_CBC_PAD)
	// El siguiente es un pseudocódigo/ejemplo conceptual usando CBC (menos ideal que GCM).
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)
	if err := m.ctx.EncryptInit(m.session, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		return nil, err
	}
	// EncryptUpdate / EncryptFinal
	out, err := m.ctx.Encrypt(m.session, plaintext)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *Manager) DecryptWithAESKey(keyHandle pkcs11.ObjectHandle, iv, ciphertext []byte) ([]byte, error) {
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, iv)
	if err := m.ctx.DecryptInit(m.session, []*pkcs11.Mechanism{mech}, keyHandle); err != nil {
		return nil, err
	}
	out, err := m.ctx.Decrypt(m.session, ciphertext)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FindKeyByLabel returns the object handle for a given label, or 0 if not found
func (m *Manager) FindKeyByLabel(label string) pkcs11.ObjectHandle {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	// Inicializar búsqueda
	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		return 0
	}
	defer m.ctx.FindObjectsFinal(m.session)

	// Buscar objetos (máximo 1)
	handles, _, err := m.ctx.FindObjects(m.session, 1)
	if err != nil {
		return 0
	}
	if len(handles) == 0 {
		return 0
	}
	return handles[0]
}

// GetAESKeyHandleByLabel returns the AES key handle for a given label, or an error if not found
func (m *Manager) GetAESKeyHandleByLabel(label string) (pkcs11.ObjectHandle, error) {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
	}

	// Inicializar búsqueda
	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		return 0, err
	}
	defer m.ctx.FindObjectsFinal(m.session)

	// Buscar objetos (máximo 1)
	handles, _, err := m.ctx.FindObjects(m.session, 1)
	if err != nil {
		return 0, err
	}
	if len(handles) == 0 {
		return 0, errors.New("AES key not found")
	}
	return handles[0], nil
}
