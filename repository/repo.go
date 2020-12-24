package repository

import (
	"fmt"
	"os"
	"path"

	"github.com/neocortical/got/index"
	"github.com/neocortical/got/object"
	"github.com/neocortical/got/ref"
)

const (
	// GitDir is the directory under which all git state is stored.
	GitDir        = ".git"
	indexFilename = "index"
	databaseDir   = "objects"
	refsDir       = "refs"
)

type Repo struct {
	workspaceDir string
	idx          index.Index
	db           object.Database
	refs         ref.Refs
}

func NewRepo(workspaceDir string) *Repo {
	return &Repo{
		workspaceDir: workspaceDir,
	}
}

func Init(workspaceDir string) (result *Repo, err error) {
	gitDir := path.Join(workspaceDir, GitDir)
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating '%s' directory: %w", GitDir, err)
	}
	for _, subDir := range []string{databaseDir, refsDir} {
		dir := path.Join(gitDir, subDir)
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("error creating '%s' directory: %v", subDir, err)
		}
	}

	return NewRepo(workspaceDir), nil
}

func (r *Repo) Dir() string {
	return path.Join(r.workspaceDir, GitDir)
}

func (r *Repo) Database() object.Database {
	if r.db == nil {
		r.db = object.NewDatabase(path.Join(r.workspaceDir, GitDir, databaseDir))
	}

	return r.db
}

func (r *Repo) Index() index.Index {
	if r.idx == nil {
		r.idx = index.NewIndex(path.Join(r.workspaceDir, GitDir, indexFilename))
	}

	return r.idx
}

func (r *Repo) Refs() ref.Refs {
	if r.refs == nil {
		r.refs = ref.NewRefs(path.Join(r.workspaceDir, GitDir))
	}

	return r.refs
}
