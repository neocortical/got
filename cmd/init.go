package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/neocortical/got/repository"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new repository.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  executeInit,
}

func executeInit(cmd *cobra.Command, args []string) error {
	basePath, err := getInitPath(args)
	if err != nil {
		return fmt.Errorf("error parsing workspace directory: %w", err)
	}

	repo, err := repository.Init(basePath)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Initialized empty Git repository in %s\n", repo.Dir())
	return nil
}

func getInitPath(args []string) (path string, err error) {
	path = wd

	if len(args) > 0 {
		path = args[0]
		path, err = filepath.Abs(path)
		if err != nil {
			return path, fmt.Errorf("error resolving path: %w", err)
		}
	}

	return
}
