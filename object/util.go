package object

import (
	"crypto/sha1"
	"fmt"
)

func generateOID(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}
