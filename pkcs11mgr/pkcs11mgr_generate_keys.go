package pkcs11mgr

import (
	"errors"

	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/utils"
)

// GenerateAESKey creates an AES key object inside SoftHSM and returns its object handle (as uint)
func GenerateAESKey(label string, id int32, bits int, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Generating AES key: label=%s, bits=%d", label, bits)
	mech := pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, utils.Int32ToByte(id)),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, bits/8),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true), // store persistently in token
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}
	if id == 0 {
		logger.AppLog.Info("The id is zero, return the last id + 1")
		var err error
		id, err = ReturnLastIDForLabel(label, s)
		if err != nil {
			logger.AppLog.Errorf("Error detect %s", err)
			return 0, err
		}
	}

	// Check if key already exists before creating it
	existingHandle, err := FindKey(label, id, s)
	if err == nil && existingHandle != 0 {
		logger.AppLog.Infof("Key with label '%s' already exists, returning existing handle: %v", label, existingHandle)
		return existingHandle, errors.New("the key is in the SSM")
	}

	handle, err := s.Ctx.GenerateKey(s.Handle, []*pkcs11.Mechanism{mech}, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate AES key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("AES key generated successfully: handle=%v", handle)
	return handle, nil
}

// GenerateDESKey creates an DES key object inside SoftHSM and returns its object handle (as uint)
func GenerateDESKey(label string, id int32, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Generating DES key: label=%s", label)

	mech := pkcs11.NewMechanism(pkcs11.CKM_DES_KEY_GEN, nil)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, utils.Int32ToByte(id)),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_DES),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}
	if id == 0 {
		logger.AppLog.Info("The id is zero, return rhe last id + 1")
		var err error
		id, err = ReturnLastIDForLabel(label, s)
		if err != nil {
			logger.AppLog.Errorf("Error detect %s", err)
			return 0, err
		}
	}
	// Check if key already exists before creating it
	existingHandle, err := FindKey(label, id, s)
	if err == nil && existingHandle != 0 {
		logger.AppLog.Infof("Key with label '%s' already exists, returning existing handle: %v", label, existingHandle)
		return existingHandle, errors.New("the key is in the SSM")
	}

	handle, err := s.Ctx.GenerateKey(s.Handle, []*pkcs11.Mechanism{mech}, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate DES key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("DES key generated successfully: handle=%v", handle)
	return handle, nil
}

// GenerateDES3Key creates an DES3 key object inside SoftHSM and returns its object handle (as uint)
func GenerateDES3Key(label string, id int32, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Generating DES3 key: label=%s", label)

	mech := pkcs11.NewMechanism(pkcs11.CKM_DES3_KEY_GEN, nil)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, utils.Int32ToByte(id)),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_DES3),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, true),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true), // store persistently in token
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}
	if id == 0 {
		logger.AppLog.Info("The id is zero, return rhe last id + 1")
		var err error
		id, err = ReturnLastIDForLabel(label, s)
		if err != nil {
			logger.AppLog.Errorf("Error detect %s", err)
			return 0, err
		}
	}
	// Check if key already exists before creating it
	existingHandle, err := FindKey(label, id, s)
	if err == nil && existingHandle != 0 {
		logger.AppLog.Infof("Key with label '%s' already exists, returning existing handle: %v", label, existingHandle)
		return existingHandle, errors.New("the key is in the SSM")
	}

	handle, err := s.Ctx.GenerateKey(s.Handle, []*pkcs11.Mechanism{mech}, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to generate DES3 key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("DES3 key generated successfully: handle=%v", handle)
	return handle, nil
}
