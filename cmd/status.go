package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

	err = idx.LoadForUpdate()
	if err != nil {
		return fmt.Errorf("error loading index: %w", err)
	}

	var untracked []string
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
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking workspace: %w", err)
	}

	for _, path := range untracked {
		fmt.Fprintln(stdout, "??", path)
	}

	return nil
}
