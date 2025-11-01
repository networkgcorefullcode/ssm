package pkcs11mgr

import (
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
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.AuditKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("Failed to initialize audit key search: %v", err)
		return err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	obj, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("Failed to find audit private key: %v", err)
		return err
	}

	if len(obj) == 0 {
		logger.AppLog.Warn("Audit private key not found, will generate new key pair")
		return generateAuditKeyPair(s)
	}

	auditPrivateKey = obj[0]

	// Load the corresponding public key
	pubTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.AuditKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, pubTemplate); err != nil {
		logger.AppLog.Errorf("Failed to initialize audit public key search: %v", err)
		return err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	pubObj, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("Failed to find audit public key: %v", err)
		return err
	}

	if len(pubObj) > 0 {
		auditPublicKey = pubObj[0]
		logger.AppLog.Info("Audit key pair loaded successfully")
	} else {
		logger.AppLog.Warn("Audit public key not found")
	}

	return nil
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
