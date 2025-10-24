package safe

import (
	"crypto/rand"
	"io"
	"runtime"

	"github.com/networkgcorefullcode/ssm/logger"
)

// Zero overwrites the slice with zeros and prevents reordering
func Zero(b []byte) {
	if b == nil {
		return
	}
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(b)
}

// Mlock attempts to lock memory to prevent swapping.
// On Windows uses VirtualLock, on Linux uses mlock.
// func Mlock(b []byte) error {
// 	if len(b) == 0 {
// 		return nil
// 	}

// 	// For Linux/Unix
// 	addr := uintptr(unsafe.Pointer(&b[0]))
// 	size := uintptr(len(b))
// 	const SYS_MLOCK = 149 // Linux/amd64
// 	_, _, err := syscall.Syscall(uintptr(SYS_MLOCK), addr, size, 0)
// 	if err != 0 {
// 		return err
// 	}
// 	return nil
// }

// func Munlock(b []byte) {
// 	if len(b) == 0 {
// 		return
// 	}

// 	// For Linux/Unix
// 	addr := uintptr(unsafe.Pointer(&b[0]))
// 	size := uintptr(len(b))
// 	const SYS_MUNLOCK = 150 // Linux/amd64
// 	syscall.Syscall(uintptr(SYS_MUNLOCK), addr, size, 0)
// }

// RandRead fills the slice with random bytes using the system's cryptographic generator.
func RandRead(b []byte) error {
	if len(b) == 0 {
		logger.AppLog.Debug("RandRead: empty slice, nothing to fill")
		return nil
	}

	logger.AppLog.Debugf("RandRead: generating %d random bytes", len(b))

	// Use crypto/rand which is the secure and standard way
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
