package pkcs11mgr

import (
	"errors"

	constants "github.com/networkgcorefullcode/ssm/const"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
)

var auditPrivateKey, auditPublicKey pkcs11.ObjectHandle

// GetAuditPrivateKey returns the audit private key handle
func GetAuditPrivateKey() pkcs11.ObjectHandle {
	return auditPrivateKey
}

// GetAuditPublicKey returns the audit public key handle
func GetAuditPublicKey() pkcs11.ObjectHandle {
	return auditPublicKey
}

// InitAuditKey initializes the audit private key by finding it in the HSM
func InitAuditKey(s *Session) error {
	// Try to find the private key using the utility function
	privateKeyHandle, err := findPrivateKeyByLabel(constants.AuditKeyLabel, *s)
	if err != nil {
		logger.AppLog.Warn("Audit private key not found, will generate new key pair")
		return generateAuditKeyPair(s)
	}

	auditPrivateKey = privateKeyHandle

	// Try to find the corresponding public key
	publicKeyHandle, err := findPublicKeyByLabel(constants.AuditKeyLabel, *s)
	if err != nil {
		logger.AppLog.Warn("Audit public key not found")
	} else {
		auditPublicKey = publicKeyHandle
		logger.AppLog.Info("Audit key pair loaded successfully")
	}

	return nil
}

// findPrivateKeyByLabel finds a private key by its label
func findPrivateKeyByLabel(label string, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for private key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return 0, err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	handles, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed: %v", err)
		return 0, err
	}
	if len(handles) == 0 {
		logger.AppLog.Warnf("No private key found with label: %s", label)
		return 0, errors.New("error Private Key With The Label Not Found")
	}
	logger.AppLog.Infof("Private key found: handle=%v", handles[0])
	return handles[0], nil
}

// findPublicKeyByLabel finds a public key by its label
func findPublicKeyByLabel(label string, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for public key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return 0, err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	handles, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed: %v", err)
		return 0, err
	}
	if len(handles) == 0 {
		logger.AppLog.Warnf("No public key found with label: %s", label)
		return 0, errors.New("error Public Key With The Label Not Found")
	}
	logger.AppLog.Infof("Public key found: handle=%v", handles[0])
	return handles[0], nil
}

// generateAuditKeyPair generates a new RSA key pair for audit signing
func generateAuditKeyPair(s *Session) error {
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.AuditKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 1}),
	}

	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.AuditKeyLabel),
	}

	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)}

	pubKey, privKey, err := s.Ctx.GenerateKeyPair(
		s.Handle,
		mechanism,
		publicKeyTemplate,
		privateKeyTemplate,
	)

	if err != nil {
		logger.AppLog.Errorf("Failed to generate audit key pair: %v", err)
		return err
	}

	auditPrivateKey = privKey
	auditPublicKey = pubKey
	logger.AppLog.Infof("Generated new audit key pair - Public: %d, Private: %d", pubKey, privKey)
	return nil
}
