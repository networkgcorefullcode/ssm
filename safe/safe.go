package safe

import (
	"runtime"
	"syscall"
	"unsafe"
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
