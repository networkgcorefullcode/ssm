package pkcs11mgr

var mgr *Manager

func SetPKCS11Manager(manager *Manager) {
	mgr = manager
}

func init() {
	session := mgr.GetSession()
	defer mgr.LogoutSession(session)

	InitAuditKey(session)
}
