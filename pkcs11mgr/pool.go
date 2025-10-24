package pkcs11mgr

// import (
// 	"errors"
// 	"sync"
// 	"time"

// 	"github.com/networkgcorefullcode/ssm/logger"
// )

// // ConnectionPool manages a pool of PKCS11 connections for reuse
// type ConnectionPool struct {
// 	pool        chan *Manager
// 	maxSize     int
// 	currentSize int
// 	pkcsPath    string
// 	slotNumber  uint
// 	pin         string
// 	mutex       sync.RWMutex
// 	closed      bool
// 	idleTimeout time.Duration
// 	maxLifetime time.Duration
// }

// // PoolConfig holds configuration for the connection pool
// type PoolConfig struct {
// 	MaxSize     int           `yaml:"maxSize"`     // Maximum number of connections in pool
// 	MinSize     int           `yaml:"minSize"`     // Minimum number of connections to maintain
// 	IdleTimeout time.Duration `yaml:"idleTimeout"` // How long a connection can be idle before being closed
// 	MaxLifetime time.Duration `yaml:"maxLifetime"` // Maximum lifetime of a connection
// 	PkcsPath    string        `yaml:"pkcsPath"`    // Path to PKCS11 library
// 	SlotNumber  uint          `yaml:"slotNumber"`  // HSM slot number
// 	Pin         string        `yaml:"pin"`         // HSM PIN
// }

// // DefaultPoolConfig returns a default pool configuration
// func DefaultPoolConfig() *PoolConfig {
// 	return &PoolConfig{
// 		MaxSize:     10,
// 		MinSize:     2,
// 		IdleTimeout: 5 * time.Minute,
// 		MaxLifetime: 30 * time.Minute,
// 	}
// }

// // Global pool instance
// var (
// 	globalPool *ConnectionPool
// 	poolOnce   sync.Once
// )

// // NewConnectionPool creates a new PKCS11 connection pool
// func NewConnectionPool(config *PoolConfig) (*ConnectionPool, error) {
// 	if config == nil {
// 		config = DefaultPoolConfig()
// 	}

// 	pool := &ConnectionPool{
// 		pool:        make(chan *Manager, config.MaxSize),
// 		maxSize:     config.MaxSize,
// 		pkcsPath:    config.PkcsPath,
// 		slotNumber:  config.SlotNumber,
// 		pin:         config.Pin,
// 		idleTimeout: config.IdleTimeout,
// 		maxLifetime: config.MaxLifetime,
// 	}

// 	// Initialize minimum connections
// 	for i := 0; i < config.MinSize; i++ {
// 		mgr, err := pool.createConnection()
// 		if err != nil {
// 			logger.AppLog.Errorf("Failed to create initial pool connection %d: %v", i, err)
// 			continue
// 		}
// 		pool.pool <- mgr
// 		pool.currentSize++
// 	}

// 	logger.AppLog.Infof("PKCS11 connection pool initialized with %d/%d connections", pool.currentSize, config.MaxSize)
// 	return pool, nil
// }

// // InitializeGlobalPool initializes the global PKCS11 connection pool
// func InitializeGlobalPool(config *PoolConfig) error {
// 	var err error
// 	poolOnce.Do(func() {
// 		globalPool, err = NewConnectionPool(config)
// 	})
// 	return err
// }

// // GetGlobalPool returns the global connection pool instance
// func GetGlobalPool() *ConnectionPool {
// 	return globalPool
// }

// // createConnection creates a new PKCS11 manager with an open session
// func (p *ConnectionPool) createConnection() (*Manager, error) {
// 	mgr, err := New(p.pkcsPath, p.slotNumber, p.pin)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = mgr.OpenSession()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set creation time for lifetime tracking
// 	mgr.createdAt = time.Now()
// 	return mgr, nil
// }

// // Get retrieves a connection from the pool or creates a new one
// func (p *ConnectionPool) Get() (*Manager, error) {
// 	p.mutex.RLock()
// 	if p.closed {
// 		p.mutex.RUnlock()
// 		return nil, errors.New("connection pool is closed")
// 	}
// 	p.mutex.RUnlock()

// 	select {
// 	case mgr := <-p.pool:
// 		// Check if connection is still valid
// 		if p.isConnectionValid(mgr) {
// 			return mgr, nil
// 		}
// 		// Connection is invalid, close it and create a new one
// 		p.closeConnection(mgr)
// 		p.mutex.Lock()
// 		p.currentSize--
// 		p.mutex.Unlock()
// 		return p.createConnection()

// 	default:
// 		// Pool is empty, try to create a new connection
// 		p.mutex.Lock()
// 		if p.currentSize < p.maxSize {
// 			p.currentSize++
// 			p.mutex.Unlock()
// 			return p.createConnection()
// 		}
// 		p.mutex.Unlock()

// 		// Pool is full, wait for a connection to be returned
// 		select {
// 		case mgr := <-p.pool:
// 			if p.isConnectionValid(mgr) {
// 				return mgr, nil
// 			}
// 			p.closeConnection(mgr)
// 			p.mutex.Lock()
// 			p.currentSize--
// 			p.mutex.Unlock()
// 			return p.createConnection()
// 		case <-time.After(30 * time.Second):
// 			return nil, errors.New("timeout waiting for available connection")
// 		}
// 	}
// }

// // Put returns a connection to the pool
// func (p *ConnectionPool) Put(mgr *Manager) {
// 	if mgr == nil {
// 		return
// 	}

// 	p.mutex.RLock()
// 	if p.closed {
// 		p.mutex.RUnlock()
// 		p.closeConnection(mgr)
// 		return
// 	}
// 	p.mutex.RUnlock()

// 	// Check if connection is still valid
// 	if !p.isConnectionValid(mgr) {
// 		p.closeConnection(mgr)
// 		p.mutex.Lock()
// 		p.currentSize--
// 		p.mutex.Unlock()
// 		return
// 	}

// 	// Update last used time
// 	mgr.lastUsed = time.Now()

// 	// Try to return to pool
// 	select {
// 	case p.pool <- mgr:
// 		// Successfully returned to pool
// 	default:
// 		// Pool is full, close the connection
// 		p.closeConnection(mgr)
// 		p.mutex.Lock()
// 		p.currentSize--
// 		p.mutex.Unlock()
// 	}
// }

// // isConnectionValid checks if a connection is still valid and not expired
// func (p *ConnectionPool) isConnectionValid(mgr *Manager) bool {
// 	if mgr == nil {
// 		return false
// 	}

// 	now := time.Now()

// 	// Check maximum lifetime
// 	if p.maxLifetime > 0 && now.Sub(mgr.createdAt) > p.maxLifetime {
// 		return false
// 	}

// 	// Check idle timeout
// 	if p.idleTimeout > 0 && now.Sub(mgr.lastUsed) > p.idleTimeout {
// 		return false
// 	}

// 	// TODO: Add actual PKCS11 session validation if needed
// 	return true
// }

// // closeConnection safely closes a PKCS11 connection
// func (p *ConnectionPool) closeConnection(mgr *Manager) {
// 	if mgr != nil {
// 		mgr.CloseSession()
// 		mgr.Finalize()
// 	}
// }

// // Close closes all connections in the pool
// func (p *ConnectionPool) Close() error {
// 	p.mutex.Lock()
// 	defer p.mutex.Unlock()

// 	if p.closed {
// 		return nil
// 	}

// 	p.closed = true
// 	close(p.pool)

// 	// Close all connections in the pool
// 	for mgr := range p.pool {
// 		p.closeConnection(mgr)
// 	}

// 	logger.AppLog.Info("PKCS11 connection pool closed")
// 	return nil
// }

// // Stats returns pool statistics
// func (p *ConnectionPool) Stats() PoolStats {
// 	p.mutex.RLock()
// 	defer p.mutex.RUnlock()

// 	return PoolStats{
// 		MaxSize:     p.maxSize,
// 		CurrentSize: p.currentSize,
// 		Available:   len(p.pool),
// 		InUse:       p.currentSize - len(p.pool),
// 	}
// }

// // PoolStats contains pool statistics
// type PoolStats struct {
// 	MaxSize     int `json:"max_size"`
// 	CurrentSize int `json:"current_size"`
// 	Available   int `json:"available"`
// 	InUse       int `json:"in_use"`
// }

// // WithConnection executes a function with a pooled connection
// func (p *ConnectionPool) WithConnection(fn func(*Manager) error) error {
// 	mgr, err := p.Get()
// 	if err != nil {
// 		return err
// 	}
// 	defer p.Put(mgr)

// 	return fn(mgr)
// }

// // Global convenience functions
// func WithConnection(fn func(*Manager) error) error {
// 	if globalPool == nil {
// 		return errors.New("global pool not initialized")
// 	}
// 	return globalPool.WithConnection(fn)
// }

// func GetPoolStats() PoolStats {
// 	if globalPool == nil {
// 		return PoolStats{}
// 	}
// 	return globalPool.Stats()
// }
