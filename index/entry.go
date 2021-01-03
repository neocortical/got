package index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

type entryHeader struct {
	CtimeSec  uint32
	CtimeNsec uint32
	MtimeSec  uint32
	MtimeNsec uint32
	Device    uint32
	INode     uint32
	Mode      uint32
	UID       uint32
	GID       uint32
	Size      uint32
	OID       [20]byte
	Flags     uint16
}

type Entry struct {
	header   entryHeader
	name     string
	pathname string
	oid      string
}

func NewEntry(pathname string, oid string, stat os.FileInfo) *Entry {
	// TODO: system-dependent filestat info
	statT := stat.Sys().(*syscall.Stat_t)

	var mode = entryModeRegular
	if stat.Mode().Perm()&0100 != 0 {
		mode = entryModeExecutable
	}

	var pathlength = len(pathname)
	if pathlength > maxPathSize {
		pathlength = maxPathSize
	}

	oidBytes, _ := hex.DecodeString(oid)

	header := entryHeader{
		CtimeSec:  uint32(statT.Ctimespec.Sec),
		CtimeNsec: uint32(statT.Ctimespec.Nsec),
		MtimeSec:  uint32(statT.Mtimespec.Sec),
		MtimeNsec: uint32(statT.Mtimespec.Nsec),
		Device:    uint32(statT.Dev),
		INode:     uint32(statT.Ino),
		Mode:      uint32(mode),
		UID:       uint32(statT.Uid),
		GID:       uint32(statT.Gid),
		Size:      uint32(statT.Size),
		Flags:     uint16(pathlength),
	}
	copy(header.OID[:], oidBytes)

	return &Entry{
		header:   header,
		name:     filepath.Base(pathname),
		pathname: pathname,
		oid:      oid,
	}
}

func (e *Entry) ParentDirectories() (result []string) {
	return parentDirectoriesForPath(e.pathname)
}

func parentDirectoriesForPath(p string) (result []string) {
	for p != "" {
		p, _ = path.Split(p)
		if p != "" {
			if p[len(p)-1] == filepath.Separator {
				p = p[:len(p)-1]
			}
			result = append([]string{p}, result...)
		}
	}

	return
}

func (e *Entry) Basename() string {
	return filepath.Base(e.pathname)
}

func (e *Entry) Path() string {
	return e.pathname
}

func (e *Entry) Encode() []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, e.header)
	if err != nil {
		fmt.Println("error writing entry header:", err)
	}

	var pathBytes = []byte(e.pathname)
	pathBytes = append(pathBytes, '\x00')
	buf.Write(pathBytes)

	for len(buf.Bytes())%8 != 0 {
		buf.WriteRune('\x00')
	}

	return buf.Bytes()
}

func (e *Entry) ModeString() string {
	return fmt.Sprintf("%o", e.header.Mode)
}

func (e *Entry) Name() string {
	return e.Basename()
}

func (e *Entry) OID() string {
	return e.oid
}
