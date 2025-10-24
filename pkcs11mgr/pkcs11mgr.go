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
	ctx        *pkcs11.Ctx
	slot       uint
	pin        string
	createdAt  time.Time
	lastUsed   time.Time
	isLoggedIn bool       // ✅ Track login state
	loginMutex sync.Mutex // ✅ Protect login operations
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
		ctx:        ctx,
		slot:       slot,
		pin:        pin,
		createdAt:  now,
		lastUsed:   now,
		isLoggedIn: false,
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
	mutexPool.Lock()
	if currentSessions >= maxSessions {
		mutexPool.Unlock()
		logger.AppLog.Debugln("The sessions reached max limit")
		return
	}
	mutexPool.Unlock()

	logger.AppLog.Debugln("Creating new PKCS#11 session")

	handle, err := m.ctx.OpenSession(m.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		logger.AppLog.Errorf("Failed to open session: %v", err)
		return
	}

	// ✅ Login solo una vez por slot (thread-safe)
	m.loginMutex.Lock()
	if !m.isLoggedIn {
		logger.AppLog.Debugf("Logging into slot %d with session: %d", m.slot, handle)
		if err := m.ctx.Login(handle, pkcs11.CKU_USER, m.pin); err != nil {
			// ✅ Verificar si el error es porque ya está logueado
			if err.Error() == "pkcs11: 0x100: CKR_USER_ALREADY_LOGGED_IN" {
				logger.AppLog.Debugf("Slot %d already logged in, continuing", m.slot)
				m.isLoggedIn = true
			} else {
				// ❌ Error real de login
				logger.AppLog.Errorf("Failed to login to session: %v", err)
				m.loginMutex.Unlock()
				_ = m.ctx.CloseSession(handle)
				return
			}
		} else {
			m.isLoggedIn = true
			logger.AppLog.Debugf("Successfully logged into slot %d", m.slot)
		}
	}
	m.loginMutex.Unlock()

	m.lastUsed = time.Now()
	session := &Session{
		Handle: handle,
		Ctx:    m.ctx,
	}

	logger.AppLog.Debugf("PKCS#11 session created: %d (total sessions: %d)", handle, currentSessions+1)

	mutexPool.Lock()
	currentSessions++
	mutexPool.Unlock()

	SessionPool <- session
}

func (m *Manager) GetSession() *Session {
	// ✅ No usar lock aquí para evitar deadlock
	if len(SessionPool) == 0 {
		mutexPool.Lock()
		needsNewSession := currentSessions < maxSessions
		mutexPool.Unlock()

		if needsNewSession {
			m.NewSession()
		}
	}

	session := <-SessionPool
	logger.AppLog.Debugf("Got session from pool: %d", session.Handle)
	return session
}

// LogoutSession returns session to pool (NO hace logout)
func (m *Manager) LogoutSession(session *Session) {
	if session == nil || session.Handle == 0 {
		return
	}
	logger.AppLog.Debugf("Returning PKCS#11 session to pool: %d", session.Handle)
	SessionPool <- session
}

// CloseSession closes a session definitively
func (m *Manager) CloseSession(session *Session) {
	if session == nil || session.Handle == 0 {
		return
	}

	logger.AppLog.Debugf("Closing PKCS#11 session: %d", session.Handle)

	// ✅ Solo cerrar la sesión, NO hacer logout (afectaría otras sesiones)
	_ = m.ctx.CloseSession(session.Handle)

	mutexPool.Lock()
	currentSessions--
	logger.AppLog.Debugf("Session closed. Remaining sessions: %d", currentSessions)
	mutexPool.Unlock()
}

// CloseAllSessions closes all sessions and does logout
func (m *Manager) CloseAllSessions() {
	m.loginMutex.Lock()
	defer m.loginMutex.Unlock()

	logger.AppLog.Debugf("Closing all sessions for slot %d", m.slot)

	// ✅ Hacer logout antes de cerrar todas las sesiones
	if m.isLoggedIn {
		// Necesitamos una sesión válida para hacer logout
		if len(SessionPool) > 0 {
			session := <-SessionPool
			_ = m.ctx.Logout(session.Handle)
			SessionPool <- session
		}
		m.isLoggedIn = false
	}

	_ = m.ctx.CloseAllSessions(m.slot)

	mutexPool.Lock()
	currentSessions = 0
	mutexPool.Unlock()
}

// Finalize cleans up the PKCS#11 context
func (m *Manager) Finalize() {
	if m.ctx != nil {
		logger.AppLog.Infoln("Finalizing PKCS#11 context")

		// ✅ Cerrar todas las sesiones antes de finalizar
		m.CloseAllSessions()

		_ = m.ctx.Finalize()
		m.ctx.Destroy()
		m.ctx = nil
	}
}
