package utils

import "encoding/binary"

// Int32ToByte converts an int32 to a big-endian 4-byte CKA_ID attribute.
// Returns nil if id is zero.
func Int32ToByte(id int32) []byte {
	if id == 0 {
		return nil
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(id))
	return b
}

// ByteToInt32 converts a big-endian 4-byte slice to an int32.
// Returns 0 if the byte slice is empty or invalid.
func ByteToInt32(b []byte) int32 {
	if len(b) == 0 {
		return 0
	}
	if len(b) < 4 {
		// Handle shorter byte slices by padding with zeros
		padded := make([]byte, 4)
		copy(padded[4-len(b):], b)
		b = padded
	}
	return int32(binary.BigEndian.Uint32(b))
}
