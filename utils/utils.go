package utils

import (
	"io"
	"crypto/rand"
	"fmt"
)

// GenerateUUID returns a UUID as a string based on RFC 4122
func GenerateUUID() string {
	uuid := GenerateBytesUUID()
	return uuidBytesToStr(uuid)
}

// GenerateBytesUUID returns a UUID as []byte based on RFC 4122
func GenerateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err)
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}

func uuidBytesToStr(uuid []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}