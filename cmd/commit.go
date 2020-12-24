package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/neocortical/got/ref"
	"github.com/neocortical/got/repository"
	"github.com/neocortical/got/tree"
	"github.com/spf13/cobra"
)

const (
	EnvAuthorName  = "GIT_AUTHOR_NAME"
	EnvAuthorEmail = "GIT_AUTHOR_EMAIL"
)

var (
	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Commit staged changes to the repository.",
		RunE:  executeCommit,
	}
	commitMessage string
)

func init() {
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
}

func executeCommit(cmd *cobra.Command, args []string) (err error) {
	workspaceDir := wd

	repo := repository.NewRepo(workspaceDir)
	db := repo.Database()
	idx := repo.Index()
	refs := repo.Refs()

	err = idx.Load()
	if err != nil {
		return fmt.Errorf("error loading index: %w", err)
	}

	t, err := tree.BuildFromIndex(idx.Entries())
	if err != nil {
		return fmt.Errorf("error building tree: %w", err)
	}

	err = t.Traverse(func(tr *tree.Tree) (oid string, err error) {
		oid, err = db.Store(tr)
		return
	})
	if err != nil {
		return fmt.Errorf("error storing tree: %w", err)
	}

	parentCommit, err := refs.ReadHead()
	if err != nil {
		return fmt.Errorf("error reading head: %w", err)
	}

	commit := ref.NewCommit(parentCommit, t.OID(), ref.Author{Name: getenv(EnvAuthorName), Email: getenv(EnvAuthorEmail), Time: time.Now()}, commitMessage)
	commitOID, err := db.Store(commit)
	if err != nil {
		return fmt.Errorf("error storing commit: %w", err)
	}

	err = refs.UpdateHead(commitOID)
	if err != nil {
		return fmt.Errorf("error storing commit SHA at HEAD: %w", err)
	}

	if parentCommit == "" {
		parentCommit = "(root-commit)"
	}

	messageStub := truncateCommitMessage(commitMessage)
	fmt.Fprintf(stdout, "[%s %s] %s\n", parentCommit, commitOID, messageStub)

	return nil
}

func truncateCommitMessage(msg string) string {
	newlineIdx := strings.Index(msg, "\n")
	if newlineIdx == -1 {
		return msg
	}

	return msg[:newlineIdx]
}
