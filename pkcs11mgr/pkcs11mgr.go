package pkcs11mgr

import (
	"errors"
	"sync"
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
var SessionPool chan *Session
var maxSessions int
var currentSessions int = 0
var mutexPool sync.Mutex

func SetChanMaxSessions(n int) {
	SessionPool = make(chan *Session, n)
	maxSessions = n
}

func (m *Manager) NewSession() {
	if currentSessions >= maxSessions {
		logger.AppLog.Debugln("The sessions get the max interval")
		return
	}
	logger.AppLog.Debugln("Creating new PKCS#11 session")

	handle, err := m.ctx.OpenSession(m.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		logger.AppLog.Errorf("Failed to open session: %v", err)
		return
	}

	// Login una sola vez al crear la sesión
	logger.AppLog.Debugf("Logging into session: %d", handle)
	if err := m.ctx.Login(handle, pkcs11.CKU_USER, m.pin); err != nil {
		logger.AppLog.Errorf("Failed to login to session: %v", err)
		_ = m.ctx.CloseSession(handle)
		return
	}

	m.lastUsed = time.Now()
	session := &Session{
		Handle: handle,
		Ctx:    m.ctx,
	}

	logger.AppLog.Debugf("PKCS#11 session created and logged in: %d", handle)
	currentSessions++
	SessionPool <- session
}

// func (m *Manager) Login(s *Session) error {
// 	if err := m.ctx.Login(s.Handle, pkcs11.CKU_USER, m.pin); err != nil {
// 		_ = m.ctx.CloseSession(s.Handle)
// 		logger.AppLog.Errorf("Failed to login to session: %v", err)
// 		return err
// 	}
// 	return nil
// }

func (m *Manager) GetSession() *Session {
	mutexPool.Lock()
	defer mutexPool.Unlock()
	if len(SessionPool) == 0 && currentSessions < maxSessions {
		m.NewSession()
	}
	session := <-SessionPool
	return session
}

// CloseSession closes an independent session
func (m *Manager) LogoutSession(session *Session) {
	if session == nil || session.Handle == 0 {
		return
	}
	logger.AppLog.Debugf("Returning PKCS#11 session to pool: %d", session.Handle)
	// ❌ NO hacer logout - mantener sesión logueada para reutilizar
	SessionPool <- session
}

// Nueva función para cerrar sesión definitivamente
func (m *Manager) CloseSession(session *Session) {
	if session == nil || session.Handle == 0 {
		return
	}
	logger.AppLog.Debugf("Closing PKCS#11 session: %d", session.Handle)
	_ = m.ctx.Logout(session.Handle)
	_ = m.ctx.CloseSession(session.Handle)

	mutexPool.Lock()
	currentSessions--
	mutexPool.Unlock()
}

// CloseSession closes an independent session
func (m *Manager) CloseAllSessions() {
	logger.AppLog.Debugf("Close all sessions for the slot %d", m.slot)
	_ = m.ctx.CloseAllSessions(m.slot)
}

// func (m *Manager) CloseSession() {
// 	for v := range SessionPool {
// 		_ = m.ctx.CloseSession(v.Handle)
// 		logger.AppLog.Debugf("Logout PKCS#11 session: %d", v.Handle)
// 	}
// }

// // GetSessionHandle returns the session handle (for compatibility with existing code)
// func (s *Session) GetHandle() pkcs11.SessionHandle {
// 	return s.Handle
// }

// Finalize cleans up the PKCS#11 context
func (m *Manager) Finalize() {
	if m.ctx != nil {
		logger.AppLog.Infoln("Finalizing PKCS#11 context")
		_ = m.ctx.Finalize()
		m.ctx.Destroy()
		m.ctx = nil
	}
}
