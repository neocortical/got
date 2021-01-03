package object

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/neocortical/got/lock"
)

type Storable interface {
	Type() string
	Serialize() []byte
}

type Database interface {
	Store(s Storable) (oid string, err error)
	Read(oid string) (result Storable, err error)
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
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %d\x00", s.Type(), len(data)))
	buf.Write(data)
	objData := buf.Bytes()

	oid = GenerateOID(objData)

	objectFilename := db.objectPath(oid)

	// short circuit if object exists
	if _, err = os.Stat(objectFilename); err == nil {
		return
	}

	dir, _ := filepath.Split(objectFilename)

	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModeDir|0755)
		}
	}
	if err != nil {
		return
	}

	var buf2 bytes.Buffer
	w := zlib.NewWriter(&buf2)
	_, err = w.Write(objData)
	if err != nil {
		return
	}

	w.Close()

	l := lock.NewLockfile(objectFilename)
	err = l.Acquire()
	if err != nil {
		return oid, fmt.Errorf("Unable to lock object for writing: %w", err)
	}

	err = l.Write(buf2.Bytes())
	if err != nil {
		return oid, fmt.Errorf("Unable to write object: %w", err)
	}

	err = l.Commit()
	if err != nil {
		return oid, fmt.Errorf("Error committing object to database: %w", err)
	}

	return
}

func (db *database) Read(oid string) (_ Storable, err error) {
	result := &genericStorable{}

	f, err := os.Open(db.objectPath(oid))
	if err != nil {
		return nil, err
	}

	unzipper, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(unzipper)

	objectType, err := buf.ReadString(0x20)
	if err != nil {
		return nil, err
	}
	result.storableType = objectType[:len(objectType)-1]

	sizeString, err := buf.ReadString(0x00)
	if err != nil {
		return nil, err
	}

	result.size, err = strconv.Atoi(sizeString[:len(sizeString)-1])
	if err != nil {
		return nil, err
	}

	result.data, err = ioutil.ReadAll(buf)
	return result, err
}
