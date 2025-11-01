package pkcs11mgr

import (
	"errors"

	"github.com/miekg/pkcs11"
	ssm_consts "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/utils"
)

// StoreKey creates a key object inside SoftHSM from raw key bytes and returns its object handle
func StoreKey(label string, key []byte, id int32, keyType string, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Storing key: label=%s, keyType=%s, keyLen=%d", label, keyType, len(key))
	var keyTypeuint uint
	switch keyType {
	case ssm_consts.TYPE_AES:
		keyTypeuint = pkcs11.CKK_AES
	case ssm_consts.TYPE_DES3:
		keyTypeuint = pkcs11.CKK_DES3
	case ssm_consts.TYPE_DES:
		keyTypeuint = pkcs11.CKK_DES
	default:
		return 0, errors.New("unsupported key type")
	}
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_ID, utils.Int32ToByte(id)),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, keyTypeuint),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, key),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, false),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_WRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_UNWRAP, false),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}

	// Check if key already exists before creating it
	existingHandle, err := FindKey(label, id, s)
	if err == nil && existingHandle != 0 {
		logger.AppLog.Infof("Key with label '%s' already exists, returning existing handle: %v", label, existingHandle)
		return existingHandle, errors.New("the key is in the SSM")
	}

	handle, err := s.Ctx.CreateObject(s.Handle, template)
	if err != nil {
		logger.AppLog.Errorf("Failed to store key: %v", err)
		return 0, err
	}
	logger.AppLog.Infof("Key stored successfully: handle=%v", handle)
	return handle, nil
}

// DeleteKey removes a key object from the HSM by label and optionally by ID
func DeleteKey(label string, id int32, s Session) error {
	logger.AppLog.Infof("Attempting to delete key with label: %s", label)

	// Find the key first
	handle, err := FindKey(label, id, s)
	if err != nil {
		logger.AppLog.Errorf("Failed to find key for deletion: %v", err)
		return err
	}

	// If key not found (handle is 0), return error
	if handle == 0 {
		logger.AppLog.Errorf("Key with label '%s' not found for deletion", label)
		return errors.New("key not found")
	}

	// Delete the key object
	if err := s.Ctx.DestroyObject(s.Handle, handle); err != nil {
		logger.AppLog.Errorf("Failed to delete key with handle %v: %v", handle, err)
		return err
	}

	logger.AppLog.Infof("Key successfully deleted: label=%s, handle=%v", label, handle)
	return nil
}

// DeleteAllKeys removes all secret key objects from the current s.Handle/slot
// WARNING: This will delete ALL keys in the token - use with caution!
func DeleteAllKeys(s Session) error {
	logger.AppLog.Warnln("Attempting to delete ALL keys from the HSM token")

	// Search for all secret key objects
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed for delete all: %v", err)
		return err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	// Find all matching objects (up to 1000 keys)
	handles, _, err := s.Ctx.FindObjects(s.Handle, 1000)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed for delete all: %v", err)
		return err
	}

	if len(handles) == 0 {
		logger.AppLog.Infoln("No keys found to delete")
		return nil
	}

	logger.AppLog.Infof("Found %d keys to delete", len(handles))

	// Delete each key object
	deletedCount := 0
	for _, handle := range handles {
		if err := s.Ctx.DestroyObject(s.Handle, handle); err != nil {
			logger.AppLog.Errorf("Failed to delete key with handle %v: %v", handle, err)
			// Continue with other keys even if one fails
			continue
		}
		deletedCount++
		logger.AppLog.Infof("Deleted key with handle: %v", handle)
	}

	logger.AppLog.Infof("Successfully deleted %d out of %d keys", deletedCount, len(handles))

	if deletedCount != len(handles) {
		return errors.New("not all keys were successfully deleted")
	}

	return nil
}

// UpdateKey updates an existing key by deleting the old one and creating a new one with the updated value
// This function combines delete and store operations to effectively "update" a key
func UpdateKey(label string, newKeyValue []byte, id int32, keyType string, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Updating key: label=%s, keyType=%s, newKeyLen=%d", label, keyType, len(newKeyValue))

	// First, check if the key exists
	existingHandle, err := FindKey(label, id, s)
	if err != nil {
		logger.AppLog.Errorf("Error searching for key to update: %v", err)
		return 0, err
	}

	if existingHandle == 0 {
		logger.AppLog.Errorf("Key with label '%s' not found for update", label)
		return 0, errors.New("key not found for update")
	}

	logger.AppLog.Infof("Found existing key to update: handle=%v", existingHandle)

	// Delete the existing key
	if err := DeleteKey(label, id, s); err != nil {
		logger.AppLog.Errorf("Failed to delete existing key for update: %v", err)
		return 0, err
	}

	logger.AppLog.Infof("Existing key deleted, creating new key with updated value")

	newHandle, err := StoreKey(label, newKeyValue, id, keyType, s)
	if err != nil {
		logger.AppLog.Errorf("Failed to create updated key: %v", err)
		return 0, err
	}

	logger.AppLog.Infof("Key updated successfully: label=%s, newHandle=%v", label, newHandle)
	return newHandle, nil
}
