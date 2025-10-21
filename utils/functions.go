package utils

import "encoding/binary"

// IntToByte converts an int32 to a big-endian 4-byte CKA_ID attribute.
// Returns nil if id is zero.
func Int32ToByte(id int32) []byte {
	if id == 0 {
		return nil
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(id))
	return b
}
