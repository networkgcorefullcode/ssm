package safe

import (
	"crypto/rand"
	"io"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/networkgcorefullcode/ssm/logger"
)

// Zero sobrescribe la slice con ceros y evita reordenamientos
func Zero(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(b)
}

// Mlock intenta bloquear la memoria para evitar swap.
// En Windows usa VirtualLock, en Linux usa mlock.
func Mlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	// Para Linux/Unix
	addr := uintptr(unsafe.Pointer(&b[0]))
	size := uintptr(len(b))
	const SYS_MLOCK = 149 // Linux/amd64
	_, _, err := syscall.Syscall(uintptr(SYS_MLOCK), addr, size, 0)
	if err != 0 {
		return err
	}
	return nil
}

func Munlock(b []byte) {
	if len(b) == 0 {
		return
	}

	// Para Linux/Unix
	addr := uintptr(unsafe.Pointer(&b[0]))
	size := uintptr(len(b))
	const SYS_MUNLOCK = 150 // Linux/amd64
	syscall.Syscall(uintptr(SYS_MUNLOCK), addr, size, 0)
}

// RandRead llena el slice con bytes aleatorios usando el generador criptográfico del sistema.
func RandRead(b []byte) error {
	if len(b) == 0 {
		logger.AppLog.Debug("RandRead: empty slice, nothing to fill")
		return nil
	}

	logger.AppLog.Debugf("RandRead: generating %d random bytes", len(b))

	// Usar crypto/rand que es la forma segura y estándar
	n, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		logger.AppLog.Errorf("RandRead: failed to read from crypto/rand: %v", err)
		return err
	}

	if n != len(b) {
		logger.AppLog.Errorf("RandRead: expected %d bytes, got %d", len(b), n)
		return io.ErrShortBuffer
	}

	logger.AppLog.Debugf("RandRead: successfully generated %d random bytes", n)
	return nil
}
