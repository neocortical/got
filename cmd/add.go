package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/neocortical/got/blob"
	"github.com/neocortical/got/index"
	"github.com/neocortical/got/object"
	"github.com/neocortical/got/repository"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [file1, file2, ...]",
	Short: "Add files/directories to the index.",
	RunE:  executeAdd,
}

func executeAdd(cmd *cobra.Command, args []string) (err error) {
	workspaceDir := wd

	repo := repository.NewRepo(workspaceDir)
	db := repo.Database()
	idx := repo.Index()

	err = idx.LoadForUpdate()
	if err != nil {
		return fmt.Errorf("error loading index: %w", err)
	}

	for _, filename := range args {
		fullPath := toAbsolutePath(filename)
		fileInfo, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			idx.Rollback()
			return fmt.Errorf("pathspec '%s' did not match any files", filename)
		}
		if err != nil {
			idx.Rollback()
			return fmt.Errorf("unexpected error adding file: %w", err)
		}

		if fileInfo.IsDir() {
			filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					if info.Name() == repository.GitDir {
						return filepath.SkipDir
					}
					return nil
				}

				err = addToIndex(db, idx, path, info)
				if err != nil {
					idx.Rollback()
					return fmt.Errorf("error adding file '%s' to index: %w", path, err)
				}

				return nil
			})
		} else {
			err = addToIndex(db, idx, filename, fileInfo)
			if err != nil {
				idx.Rollback()
				return fmt.Errorf("error adding file '%s' to index: %w", filename, err)
			}
		}
	}

	err = idx.WriteUpdates()
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error committing the index: %w", err)
	}

	return nil
}

func addToIndex(db object.Database, idx index.Index, filename string, info os.FileInfo) (err error) {
	data, err := ioutil.ReadFile(toAbsolutePath(filename))
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error reading file '%s': %w", filename, err)
	}

	b := blob.New(data)

	oid, err := db.Store(b)
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error storing blob '%s': %w", filename, err)
	}

	idx.Add(index.NewEntry(toRelativePath(filename), oid, info))

	return nil
}
