package pkcs11mgr

import (
	"errors"
	"math/rand/v2"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/utils"
)

type ObjectAttributes struct {
	Handle int32
	Id     int32
}

// FindKey returns the object handle for a given label, or 0 if not found return a one key
func FindKey(label string, id int32, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if id != 0 {
		// Convert int32 to byte slice
		b := utils.Int32ToByte(id)
		template = append(template, pkcs11.NewAttribute(pkcs11.CKA_ID, b))
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
		logger.AppLog.Warnf("No key found with label: %s", label)
		return 0, errors.New(constants.ERROR_STRING_KEY_NOT_FOUND)
	}
	logger.AppLog.Infof("Key found: handle=%v", handles[0])
	return handles[0], nil
}

// FindKey using key label as a filter returns the object handles
func FindKeysLabel(label string, s Session) ([]pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return nil, err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	var handles []pkcs11.ObjectHandle
	for {
		new_handles, _, err := s.Ctx.FindObjects(s.Handle, 20) // return a []ObjectHandle the max size is 20
		if err != nil {
			logger.AppLog.Errorf("FindObjects failed: %v", err)
			return nil, err
		}
		if len(new_handles) == 0 {
			if len(handles) == 0 {
				logger.AppLog.Warnf("No key found with label: %s", label)
				return nil, errors.New(constants.ERROR_STRING_KEY_NOT_FOUND)
			}
			logger.AppLog.Info("Key found is finished")
			return handles, nil
		}

		handles = append(handles, new_handles...)
	}
}

// FindAllKeys returns all secret keys in the HSM grouped by label
func FindAllKeys(s Session) (map[string][]pkcs11.ObjectHandle, error) {
	logger.AppLog.Info("Searching for all keys in HSM")

	// Search for all secret keys without label filter
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return nil, err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	var allHandles []pkcs11.ObjectHandle
	for {
		newHandles, _, err := s.Ctx.FindObjects(s.Handle, 20)
		if err != nil {
			logger.AppLog.Errorf("FindObjects failed: %v", err)
			return nil, err
		}
		if len(newHandles) == 0 {
			if len(allHandles) == 0 {
				logger.AppLog.Warnf("No key found")
				return map[string][]pkcs11.ObjectHandle{}, errors.New(constants.ERROR_STRING_KEY_NOT_FOUND)
			}
			break
		}
		allHandles = append(allHandles, newHandles...)
	}

	logger.AppLog.Infof("Found %d keys in HSM", len(allHandles))

	// Group handles by label
	keysByLabel := make(map[string][]pkcs11.ObjectHandle)
	for _, handle := range allHandles {
		label, err := GetObjectLabel(handle, s)
		if err != nil {
			logger.AppLog.Warnf("Failed to get label for handle %d: %v", handle, err)
			continue
		}
		keysByLabel[label] = append(keysByLabel[label], handle)
	}

	logger.AppLog.Infof("Keys grouped into %d labels", len(keysByLabel))
	return keysByLabel, nil
}

// FindKeyLabelReturnRandom returns random object handle for a given label, or 0 if not found return a one key
func FindKeyLabelReturnRandom(label string, s Session) (pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	handles, err := FindKeysLabel(label, s)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed: %v", err)
		return 0, err
	}
	return handles[rand.Int64N(int64(len(handles)))], err
}

func GetValuesForObjects(o []pkcs11.ObjectHandle, s Session) ([]ObjectAttributes, error) {
	logger.AppLog.Info("Get attributes for handles objects")

	var result []ObjectAttributes
	for _, handle := range o {
		attr, err := GetObjectAttributes(handle, s)
		if err != nil {
			logger.AppLog.Errorf("Failed to get attributes for handle %d: %v", handle, err)
			continue
		}
		result = append(result, ObjectAttributes{
			Handle: int32(handle),
			Id:     attr.Id,
		})
	}
	return result, nil
}

// GetObjectAttributes retrieves the CKA_ID and CKA_VALUE_LEN attributes for a given object handle
func GetObjectAttributes(handle pkcs11.ObjectHandle, s Session) (ObjectAttributes, error) {
	logger.AppLog.Infof("Getting attributes for object handle: %d", handle)

	// Define the attributes we want to retrieve
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
	}

	// Get the attribute values
	attrs, err := s.Ctx.GetAttributeValue(s.Handle, handle, template)
	if err != nil {
		logger.AppLog.Errorf("GetAttributeValue failed for handle %d: %v", handle, err)
		return ObjectAttributes{}, err
	}

	var result ObjectAttributes
	result.Handle = int32(handle)

	// Parse the returned attributes
	for _, attr := range attrs {
		switch attr.Type {
		case pkcs11.CKA_ID:
			if len(attr.Value) > 0 {
				result.Id = utils.ByteToInt32(attr.Value)
			}
		}
	}

	logger.AppLog.Infof("Attributes retrieved - Handle: %d, Id: %d", result.Handle, result.Id)
	return result, nil
}

// GetObjectLabel retrieves the CKA_LABEL attribute for a given object handle
func GetObjectLabel(handle pkcs11.ObjectHandle, s Session) (string, error) {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, nil),
	}

	attrs, err := s.Ctx.GetAttributeValue(s.Handle, handle, template)
	if err != nil {
		return "", err
	}

	if len(attrs) > 0 && len(attrs[0].Value) > 0 {
		return string(attrs[0].Value), nil
	}

	return "", nil
}

// ReturnLastIDForLabel get a label and return the last id used plus one
func ReturnLastIDForLabel(label string, s Session) (int32, error) {
	keys, err := FindKeysLabel(label, s)
	if err != nil && err.Error() == constants.ERROR_STRING_KEY_NOT_FOUND {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}

	var maxID int32
	for _, key := range keys {
		attr, err := GetObjectAttributes(key, s)
		if err != nil {
			return 0, err
		}
		if attr.Id > maxID {
			maxID = attr.Id
		}
	}

	return maxID + 1, nil
}
