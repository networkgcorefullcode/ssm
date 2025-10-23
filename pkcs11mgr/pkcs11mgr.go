package pkcs11mgr

import (
	"errors"
	"time"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

type Manager struct {
	ctx       *pkcs11.Ctx
	slot      uint
	session   pkcs11.SessionHandle
	pin       string
	createdAt time.Time
	lastUsed  time.Time
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
