package middleware

import (
	"os"

	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

var mgr *pkcs11mgr.Manager
var hostname string

func SetPKCS11Manager(manager *pkcs11mgr.Manager) {
	mgr = manager
}

func init() {
	hostname, _ = os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
}
