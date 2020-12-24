package ref

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/neocortical/got/lock"
)

const (
	headFilename = "HEAD"
)

type Refs interface {
	ReadHead() (result string, err error)
	UpdateHead(val string) error
}

type refs struct {
	dir string
}

func NewRefs(dir string) Refs {
	return &refs{dir}
}

func (r *refs) ReadHead() (result string, err error) {
	headPath := path.Join(r.dir, headFilename)

	_, err = os.Stat(headPath)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(headPath)
	if err != nil {
		return
	}

	return strings.TrimSpace(string(data)), nil
}

func (r *refs) UpdateHead(oid string) (err error) {
	lf := lock.NewLockfile(path.Join(r.dir, "HEAD"))

	if err = lf.Acquire(); err != nil {
		return fmt.Errorf("could not lock HEAD for writing: %w", err)
	}

	if err = lf.Write([]byte(fmt.Sprintf("%s\n", oid))); err != nil {
		return fmt.Errorf("failed to write HEAD data: %w", err)
	}

	if err = lf.Commit(); err != nil {
		return fmt.Errorf("failed to commit HEAD data: %w", err)
	}

	return
}
