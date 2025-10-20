package pkcs11mgr

import (
	"github.com/miekg/pkcs11"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/utils"
)

type ObjectAttributes struct {
	Handle   int32
	Id       int32
	SizeBits int32
}

// FindKey returns the object handle for a given label, or 0 if not found return a one key
func (m *Manager) FindKey(label string, id int32) (pkcs11.ObjectHandle, error) {
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

	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return 0, err
	}
	defer m.ctx.FindObjectsFinal(m.session)

	handles, _, err := m.ctx.FindObjects(m.session, 1)
	if err != nil {
		logger.AppLog.Errorf("FindObjects failed: %v", err)
		return 0, err
	}
	if len(handles) == 0 {
		logger.AppLog.Warnf("No key found with label: %s", label)
		return 0, err
	}
	logger.AppLog.Infof("Key found: handle=%v", handles[0])
	return handles[0], nil
}

// FindKey using key label as a filter returns the object handles
func (m *Manager) FindKeysLabel(label string) ([]pkcs11.ObjectHandle, error) {
	logger.AppLog.Infof("Searching for key by label: %s", label)
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
	}

	if err := m.ctx.FindObjectsInit(m.session, template); err != nil {
		logger.AppLog.Errorf("FindObjectsInit failed: %v", err)
		return 0, err
	}
	defer m.ctx.FindObjectsFinal(m.session)

	var handles []pkcs11.ObjectHandle
	var iteration uint32 = 0
	for {
		new_handles, _, err := m.ctx.FindObjects(m.session, 20) // return a []ObjectHandle the max size is 20
		if len(new_handles) == 0 && iteration != 0 {
			logger.AppLog.Info("Key found is finished")
			return handles, nil
		}
		if len(new_handles) == 0 && iteration == 0 {
			logger.AppLog.Info("Key found is finished, not found any key with this label")
			return 0, nil
		}
		if err != nil {
			logger.AppLog.Errorf("FindObjects failed: %v", err)
			return 0, err
		}
		handles = append(handles, new_handles...)
		iteration++
	}
}

// FindKey using key label as a filter returns the object handles
func (m *Manager) GetValuesForObjects(o []pkcs11.ObjectHandle) ([]ObjectAttributes, error) {
	logger.AppLog.Info("Get attributes for handles objects")

	var result []ObjectAttributes
	for _, handle := range o {
		attr, err := m.GetObjectAttributes(handle)
		if err != nil {
			logger.AppLog.Errorf("Failed to get attributes for handle %d: %v", handle, err)
			continue
		}
		result = append(result, ObjectAttributes{
			Handle:   int32(handle),
			Id:       attr.Id,
			SizeBits: attr.SizeBits,
		})
	}
	return result, nil
}

// GetObjectAttributes retrieves the CKA_ID and CKA_VALUE_LEN attributes for a given object handle
func (m *Manager) GetObjectAttributes(handle pkcs11.ObjectHandle) (ObjectAttributes, error) {
	logger.AppLog.Infof("Getting attributes for object handle: %d", handle)

	// Define the attributes we want to retrieve
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, nil),
	}

	// Get the attribute values
	attrs, err := m.ctx.GetAttributeValue(m.session, handle, template)
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
		case pkcs11.CKA_VALUE_LEN:
			if len(attr.Value) > 0 {
				// CKA_VALUE_LEN returns the length in bytes, convert to bits
				valueLen := utils.ByteToInt32(attr.Value)
				result.SizeBits = valueLen * 8
			}
		}
	}

	logger.AppLog.Infof("Attributes retrieved - Handle: %d, Id: %d, SizeBits: %d", result.Handle, result.Id, result.SizeBits)
	return result, nil
}
