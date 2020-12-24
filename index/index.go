package index

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/neocortical/got/lock"
)

const (
	signature           = "DIRC"
	entryModeRegular    = 0100644
	entryModeExecutable = 0100755
	maxPathSize         = 0xfff
)

var lockConflictErrTemplate = `%v

Another git process seems to be running in this repository, e.g.
an editor opened by 'git commit'. Please make sure all processes
are terminated then try again. If it still fails, a git process
may have crashed in this repository earlier:
remove the file manually to continue.`

type header struct {
	Signature [4]byte
	Version   uint32
	Entries   uint32
}

type Index interface {
	LoadForUpdate() error
	Load() error
	Add(e *Entry)
	WriteUpdates() error
	Entries() []*Entry
	Rollback()
	IsTracked(path string) bool
	FirstUntrackedPath(path string) string
}

type index struct {
	l         *lock.Lockfile
	filename  string
	entryMap  map[string]*Entry
	parentMap map[string][]string
	changed   bool
}

func NewIndex(idxFilename string) Index {
	return &index{filename: idxFilename, entryMap: map[string]*Entry{}, parentMap: map[string][]string{}}
}

func (idx *index) LoadForUpdate() (err error) {
	l := lock.NewLockfile(idx.filename)
	err = l.Acquire()
	if lock.IsLockConflict(err) {
		return fmt.Errorf(lockConflictErrTemplate, err)
	}
	if err != nil {
		return err
	}

	idx.l = l

	err = idx.Load()
	return
}

func (idx *index) Load() (err error) {

	fileData, err := ioutil.ReadFile(idx.filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading index file: %w", err)
	}

	// verify checksum
	if len(fileData) <= 20 {
		return errors.New("invalid index file format")
	}
	if !verifySHAsMatch(generateSHA(fileData[:len(fileData)-20]), fileData[len(fileData)-20:]) {
		return fmt.Errorf("Index file checksum mismatch (%x != %x)", generateSHA(fileData[:len(fileData)-20]), fileData[len(fileData)-20:])
	}

	f := bytes.NewBuffer(fileData)

	var header header
	err = binary.Read(f, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("error reading index header: %w", err)
	}

	if string(header.Signature[:]) != signature {
		return fmt.Errorf("invalid file signature bytes: %v", header.Signature)
	}
	if header.Version != 2 {
		return fmt.Errorf("invalid index file version: %d", header.Version)
	}

	for i := 0; i < int(header.Entries); i++ {
		var entry Entry
		entry, err = readNextEntry(f)
		if err != nil {
			return
		}

		idx.entryMap[entry.pathname] = &entry
		for _, dir := range entry.ParentDirectories() {
			idx.parentMap[dir] = append(idx.parentMap[dir], entry.pathname)
		}
	}

	return nil
}

func (i *index) Add(entry *Entry) {
	if existing, exists := i.entryMap[entry.pathname]; exists && existing.oid == entry.oid {
		return
	}

	i.removeConflicts(entry)

	i.entryMap[entry.pathname] = entry
	for _, dir := range entry.ParentDirectories() {
		i.parentMap[dir] = append(i.parentMap[dir], entry.pathname)
	}

	i.changed = true
}

func (i *index) removeConflicts(entry *Entry) {
	// remove any conflicting file
	for _, dir := range entry.ParentDirectories() {
		delete(i.entryMap, dir)
	}

	// remove any files that live under conflicting directories
	for _, deadentry := range i.parentMap[entry.pathname] {
		delete(i.entryMap, deadentry)
	}
	delete(i.parentMap, entry.pathname)
}

func (i *index) Entries() (result []*Entry) {
	var entrynames []string
	for k := range i.entryMap {
		entrynames = append(entrynames, k)
	}

	sort.Strings(entrynames)
	for _, pathname := range entrynames {
		result = append(result, i.entryMap[pathname])
	}

	return
}

func (i *index) WriteUpdates() (err error) {
	if !i.changed {
		if i.l != nil {
			err = i.l.Rollback()
		}
		return
	}

	buf := new(bytes.Buffer)

	header := header{
		Signature: [4]byte{'D', 'I', 'R', 'C'},
		Version:   2,
		Entries:   uint32(len(i.entryMap)),
	}
	binary.Write(buf, binary.BigEndian, &header)

	for _, e := range i.Entries() {
		buf.Write(e.Encode())
	}

	data := buf.Bytes()
	err = i.l.Write(data)
	if err != nil {
		return
	}

	err = i.l.Write(generateSHA(data))
	if err != nil {
		return
	}

	err = i.l.Commit()
	return
}

func generateSHA(data []byte) []byte {
	hasher := sha1.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func verifySHAsMatch(sha1, sha2 []byte) bool {
	if len(sha1) != len(sha2) {
		return false
	}

	for i := range sha1 {
		if sha1[i] != sha2[i] {
			return false
		}
	}

	return true
}

func readNextBytes(file *os.File, length int) ([]byte, error) {
	bytes := make([]byte, length)

	_, err := file.Read(bytes)
	return bytes, err
}

func readNextEntry(r io.Reader) (result Entry, err error) {
	var header entryHeader
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return result, fmt.Errorf("error reading entry header: %w", err)
	}

	// fmt.Printf("ctime: %d, %d\nmtime: %d, %d\noid: %x\nflags: %d\n", header.CtimeSec, header.CtimeNsec, header.MtimeSec, header.MtimeNsec, header.OID, header.Flags)

	if header.Flags < maxPathSize {
		pathnameBytes := make([]byte, header.Flags)
		_, err = r.Read(pathnameBytes)
		if err != nil {
			return result, fmt.Errorf("error reading %d pathname bytes: %w", header.Flags, err)
		}

		result.pathname = string(pathnameBytes)
		result.header = header
		result.oid = fmt.Sprintf("%x", header.OID[:])

		// advance past nulls

		nullsToRead := ((8 - (63+header.Flags)%8) % 8) + 1
		// fmt.Println("consuming nulls", nullsToRead)
		if nullsToRead > 0 {
			var nullReader = make([]byte, calculatePathnameNullsDoRead(header.Flags))
			_, err = r.Read(nullReader)
			if err != nil {
				return result, fmt.Errorf("error consuming %d nulls: %w", nullsToRead, err)
			}
		}

		return
	}

	var path []byte
	var pathBuf = make([]byte, 8)
	for {
		_, err = r.Read(pathBuf)
		if err != nil {
			return result, fmt.Errorf("error reading next 8 bytes of pathname: %w", err)
		}

		if pathBuf[7] == 0x00 {
			for i := range pathBuf {
				if pathBuf[i] == 0x00 {
					result.pathname = string(path)
					return
				}
				path = append(path, pathBuf[i])
			}
		} else {
			path = append(path, pathBuf...)
		}
	}
}

func calculatePathnameNullsDoRead(headerFlags uint16) int {
	return ((8 - (63+int(headerFlags))%8) % 8) + 1
}

func (i *index) Rollback() {
	if i.l != nil {
		_ = i.l.Rollback()
	}
}

func (i *index) IsTracked(path string) (result bool) {
	_, result = i.entryMap[path]
	return
}

func (i *index) FirstUntrackedPath(path string) string {
	if i.IsTracked(path) {
		return ""
	}

	for _, dir := range parentDirectoriesForPath(path) {
		if _, directoryTracked := i.parentMap[dir]; !directoryTracked {
			return dir + string(filepath.Separator)
		}
	}

	return path
}
