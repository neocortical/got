package object

import (
	"crypto/sha1"
	"fmt"
)

// GenerateOID generates a SHA1 object ID for use in the database, index, refs, etc.
func GenerateOID(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
