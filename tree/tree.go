package tree

import (
	"bytes"
	"encoding/hex"
	"errors"
	"path/filepath"
	"sort"

	"github.com/neocortical/got/index"
)

const (
	dirModeString = "40000"
)

type Node interface {
	ModeString() string
	Name() string
	OID() string
}

type Tree struct {
	name    string
	oid     string
	entries map[string]Node
}

func (t *Tree) ModeString() string {
	return dirModeString
}

func (t *Tree) Name() string {
	return t.name
}

func (t *Tree) OID() string {
	return t.oid
}

func (t Tree) Type() string {
	return "tree"
}

func (t *Tree) Serialize() []byte {
	var buf bytes.Buffer

	for _, node := range t.Entries() {
		buf.WriteString(node.ModeString())
		buf.WriteRune(' ')
		buf.WriteString(node.Name())
		buf.WriteRune('\x00')
		oidBytes, _ := hex.DecodeString(node.OID())
		buf.Write(oidBytes)
	}

	return buf.Bytes()
}

func BuildFromIndex(entries []*index.Entry) (result *Tree, err error) {
	t := &Tree{
		entries: map[string]Node{},
	}
	for _, e := range entries {
		err = t.AddEntry(e.ParentDirectories(), e)
		if err != nil {
			return
		}
	}

	return t, nil
}

func (t *Tree) AddEntry(parents []string, e *index.Entry) (err error) {
	//	fmt.Println("in AddEntry", parents, t, e)
	if len(parents) == 0 {
		t.entries[e.Basename()] = e
		return
	}

	var subtree *Tree
	dir := filepath.Base(parents[0])
	node, ok := t.entries[dir]
	if !ok {
		subtree = &Tree{
			name:    dir,
			entries: map[string]Node{},
		}
		t.entries[dir] = subtree
	} else {
		subtree, ok = node.(*Tree)
		if !ok {
			return errors.New("dir/file mismatch in index")
		}
	}

	err = subtree.AddEntry(parents[1:], e)
	return
}

func (t Tree) Entries() (result []Node) {
	var names []string
	for k := range t.entries {
		names = append(names, k)
	}

	sort.Strings(names)
	for _, name := range names {
		result = append(result, t.entries[name])
	}

	return
}

func (t *Tree) Traverse(store func(*Tree) (string, error)) (err error) {
	for _, node := range t.Entries() {
		switch n := node.(type) {
		case *Tree:
			err = n.Traverse(store)
			if err != nil {
				return
			}
		}
	}

	t.oid, err = store(t)
	return
}
