package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path"

	"github.com/neocortical/got/lock"
)

type Storable interface {
	Type() string
	Serialize() []byte
}

type Database interface {
	Store(s Storable) (oid string, err error)
}

type database struct {
	dir string
}

func NewDatabase(dir string) Database {
	return &database{
		dir: dir,
	}
}

func (db *database) Store(s Storable) (oid string, err error) {
	data := s.Serialize()
	oid = GenerateOID(data)

	dir := path.Join(db.dir, oid[0:2])
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModeDir|0755)
		}
	}
	if err != nil {
		return
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err = w.Write(data)
	if err != nil {
		return
	}
	w.Close()

	// short circuit if object exists
	filename := path.Join(db.dir, oid[0:2], oid[2:])
	if _, err = os.Stat(filename); err == nil {
		return
	}

	l := lock.NewLockfile(filename)
	err = l.Acquire()
	if err != nil {
		return oid, fmt.Errorf("Unable to lock object for writing: %w", err)
	}

	err = l.Write(buf.Bytes())
	if err != nil {
		return oid, fmt.Errorf("Unable to write object: %w", err)
	}

	err = l.Commit()
	if err != nil {
		return oid, fmt.Errorf("Error committing object to database: %w", err)
	}

	return
}
