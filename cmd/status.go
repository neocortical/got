package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/neocortical/got/blob"
	"github.com/neocortical/got/index"
	"github.com/neocortical/got/repository"
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "View the status of the local repository.",
		RunE:  executeStatus,
	}
)

func executeStatus(cmd *cobra.Command, args []string) (err error) {
	workspaceDir := wd

	repo := repository.NewRepo(workspaceDir)
	idx := repo.Index()
	db := repo.Database()

	err = idx.LoadForUpdate()
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error loading index: %w", err)
	}

	var untracked []string
	var modified []string
	var untrackedSet = map[string]struct{}{}
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == repository.GitDir {
				return filepath.SkipDir
			}
			return nil
		}

		relativePath := toRelativePath(path)
		if !idx.IsTracked(relativePath) {
			path := idx.FirstUntrackedPath(relativePath)
			if _, seen := untrackedSet[path]; !seen {
				untrackedSet[path] = struct{}{}
				untracked = append(untracked, path)
			}
		} else {
			statModified, timesModified := idx.IsMetadataModified(relativePath, info)
			if statModified {
				modified = append(modified, path)
			}

			if !timesModified {
				return nil
			}

			// Light modification was inconclusive. Gotta read the file and compare the content to the index
			data, err := ioutil.ReadFile(toAbsolutePath(path))
			if err != nil {
				idx.Rollback()
				return fmt.Errorf("error reading file '%s': %w", relativePath, err)
			}

			b := blob.New(data)

			oid, err := db.Store(b)
			if err != nil {
				idx.Rollback()
				return fmt.Errorf("error storing blob '%s': %w", relativePath, err)
			}

			existingEntry, _ := idx.GetEntry(relativePath)

			if oid == existingEntry.OID() {
				if timesModified {
					idx.Add(index.NewEntry(relativePath, oid, info))
				}

				return nil
			}

			modified = append(modified, relativePath)
		}

		return nil
	})
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error walking workspace: %w", err)
	}

	for _, path := range modified {
		fmt.Fprintln(stdout, " M", path)
	}

	for _, path := range untracked {
		fmt.Fprintln(stdout, "??", path)
	}

	err = idx.WriteUpdates()
	return
}
