package pkcs11mgr

import (
	"errors"
	"time"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

// Manager manages PKCS#11 context without holding a specific session
type Manager struct {
	ctx       *pkcs11.Ctx
	slot      uint
	pin       string
	createdAt time.Time
	lastUsed  time.Time
}

// Session represents an independent PKCS#11 session
type Session struct {
	Handle pkcs11.SessionHandle
	Ctx    *pkcs11.Ctx
}

func New(modulePath string, slot uint, pin string) (*Manager, error) {
	ctx := pkcs11.New(modulePath)
	if ctx == nil {
		return nil, errors.New("pkcs11.New returned nil")
	}
	if err := ctx.Initialize(); err != nil {
		return nil, err
	}

	now := time.Now()
	mgr := &Manager{
		ctx:       ctx,
		slot:      slot,
		pin:       pin,
		createdAt: now,
		lastUsed:  now,
	}
	logger.AppLog.Infoln("PKCS#11 module initialized")
	return mgr, nil
}

// NewSession creates and returns a new independent PKCS#11 session
func (m *Manager) NewSession() (*Session, error) {
	logger.AppLog.Debugln("Creating new PKCS#11 session")

	// Open a new session with the specified slot
	handle, err := m.ctx.OpenSession(m.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		logger.AppLog.Errorf("Failed to open session: %v", err)
		return nil, err
	}

	// Login to the session
	if err := m.ctx.Login(handle, pkcs11.CKU_USER, m.pin); err != nil {
		_ = m.ctx.CloseSession(handle)
		logger.AppLog.Errorf("Failed to login to session: %v", err)
		return nil, err
	}

	m.lastUsed = time.Now()
	session := &Session{
		Handle: handle,
		Ctx:    m.ctx,
	}

	logger.AppLog.Debugf("PKCS#11 session created: %d", handle)
	return session, nil
}

// CloseSession closes an independent session
func (m *Manager) CloseSession(session *Session) {
	if session == nil || session.Handle == 0 {
		return
	}

	logger.AppLog.Debugf("Closing PKCS#11 session: %d", session.Handle)
	_ = m.ctx.Logout(session.Handle)
	_ = m.ctx.CloseSession(session.Handle)
	session.Handle = 0
	logger.AppLog.Debugln("PKCS#11 session closed")
}

// GetSessionHandle returns the session handle (for compatibility with existing code)
func (s *Session) GetHandle() pkcs11.SessionHandle {
	return s.Handle
}

// Finalize cleans up the PKCS#11 context
func (m *Manager) Finalize() {
	if m.ctx != nil {
		logger.AppLog.Infoln("Finalizing PKCS#11 context")
		_ = m.ctx.Finalize()
		m.ctx.Destroy()
		m.ctx = nil
	}
}

// Legacy methods for backward compatibility (DEPRECATED - use NewSession/CloseSession instead)

// OpenSession creates a default session (DEPRECATED: use NewSession instead)
func (m *Manager) OpenSession() (*Session, error) {
	logger.AppLog.Warnln("OpenSession is deprecated, use NewSession instead")
	return m.NewSession()
}
