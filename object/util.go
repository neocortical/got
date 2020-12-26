package object

import (
	"crypto/sha1"
	"fmt"
	"path"
)

type genericStorable struct {
	storableType string
	size         int
	data         []byte
}

func (gs *genericStorable) Type() string {
	return gs.storableType
}

func (gs *genericStorable) Serialize() []byte {
	return gs.data
}

// GenerateOID generates a SHA1 object ID for use in the database, index, refs, etc.
func GenerateOID(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func (db *database) objectPath(oid string) string {
	return path.Join(db.dir, oid[0:2], oid[2:])
}
