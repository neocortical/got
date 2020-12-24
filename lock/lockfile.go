package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrLockNotHeld = errors.New("lockfile not held by this lock, acquire first")
)

type errLockConflict struct {
	lockfile string
}

func (elc *errLockConflict) Error() string {
	return fmt.Sprintf("Unable to create '%s': File exists.", elc.lockfile)
}

func IsLockConflict(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(*errLockConflict); ok {
		return true
	}

	return false
}

type Lockfile struct {
	path string
	file *os.File
}

func NewLockfile(path string) *Lockfile {
	return &Lockfile{path, nil}
}

func (lf *Lockfile) Acquire() (err error) {
	lockfile, err := filepath.Abs(fmt.Sprintf("%s.lock", lf.path))
	if err != nil {
		return
	}

	lf.file, err = os.OpenFile(lockfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return &errLockConflict{lockfile}
		}
	}

	return
}

func (lf *Lockfile) Write(data []byte) (err error) {
	if lf.file == nil {
		return ErrLockNotHeld
	}

	_, err = lf.file.Write(data)
	return
}

func (lf *Lockfile) Commit() (err error) {
	if lf.file == nil {
		return ErrLockNotHeld
	}

	err = lf.file.Close()
	if err != nil {
		return
	}

	err = os.Rename(fmt.Sprintf("%s.lock", lf.path), lf.path)
	return
}

func (lf *Lockfile) Rollback() (err error) {
	if lf.file == nil {
		return ErrLockNotHeld
	}

	err = lf.file.Close()
	if err != nil {
		return
	}

	err = os.Remove(fmt.Sprintf("%s.lock", lf.path))
	return
}
