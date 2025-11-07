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
