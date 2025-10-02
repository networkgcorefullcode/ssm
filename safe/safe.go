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
// Requiere privilegios CAP_IPC_LOCK o root.
func Mlock(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	addr := uintptr(unsafe.Pointer(&b[0]))
	size := uintptr(len(b))
	const SYS_MLOCK = 149 // Linux/amd64
	_, _, err := syscall.SyscallN(SYS_MLOCK, addr, size, 0)
	if err != 0 {
		return err
	}
	return nil
}

func Munlock(b []byte) {
	if len(b) == 0 {
		return
	}
	addr := uintptr(unsafe.Pointer(&b[0]))
	const SYS_MUNLOCK = 150 // Linux/amd64
	size := uintptr(len(b))
	syscall.SyscallN(SYS_MUNLOCK, addr, size, 0)
}
