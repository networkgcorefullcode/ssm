package handlers

import "github.com/networkgcorefullcode/ssm/pkcs11mgr"

var mgr *pkcs11mgr.Manager

func SetPKCS11Manager(manager *pkcs11mgr.Manager) {
	mgr = manager
}
