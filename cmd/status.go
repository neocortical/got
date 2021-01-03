package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/neocortical/got/blob"
	"github.com/neocortical/got/index"
	"github.com/neocortical/got/object"
	"github.com/neocortical/got/ref"
	"github.com/neocortical/got/repository"
	"github.com/neocortical/got/tree"
	"github.com/spf13/cobra"
)

const (
	statusWorkspaceModified = 0b0001
	statusWorkspaceDeleted  = 0b0010
	statusindexModified     = 0b0100
	statusIndexDeleted      = 0b1000
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
	refs := repo.Refs()

	err = idx.LoadForUpdate()
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error loading index: %w", err)
	}

	var untracked []string
	var modified = map[string]int{}
	var untrackedSet = map[string]struct{}{}
	var workspaceFileset = map[string]struct{}{}
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
		workspaceFileset[relativePath] = struct{}{}

		if !idx.IsTracked(relativePath) {
			relativePath := idx.FirstUntrackedPath(relativePath)
			if _, seen := untrackedSet[relativePath]; !seen {
				untrackedSet[relativePath] = struct{}{}
				untracked = append(untracked, relativePath)
			}
		} else {
			statModified, timesModified := idx.IsMetadataModified(relativePath, info)
			if statModified {
				modified[path] |= statusWorkspaceModified
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

			modified[relativePath] |= statusWorkspaceModified
		}

		return nil
	})
	if err != nil {
		idx.Rollback()
		return fmt.Errorf("error walking workspace: %w", err)
	}

	for _, entry := range idx.Entries() {
		// fmt.Println(entry.Name())
		if _, stillExists := workspaceFileset[entry.Path()]; !stillExists {
			modified[entry.Path()] |= statusWorkspaceDeleted
		}
	}

	var modifiedPaths []string
	for path := range modified {
		modifiedPaths = append(modifiedPaths, path)
	}
	sort.Strings(modifiedPaths)

	// cache/HEAD changes
	headCommitOID, err := refs.ReadHead()
	if err != nil {
		return fmt.Errorf("error reading head: %w", err)
	}

	if headCommitOID != "" {
		headCommitObj, err := db.Read(headCommitOID)
		if err != nil {
			return fmt.Errorf("error reading head commit from database: %w", err)
		}

		headCommit, err := ref.DeserializeCommit(headCommitObj.Serialize())
		if err != nil {
			return fmt.Errorf("error reading/parsing head commit: %w", err)
		}

		err = showTree(db, headCommit.TreeOID, "")
		if err != nil {
			return fmt.Errorf("error showing HEAD commit tree")
		}
	}

	for _, path := range modifiedPaths {
		fmt.Fprintf(stdout, "%s %s\n", porcelainStatus(modified[path]), path)
	}

	for _, path := range untracked {
		fmt.Fprintln(stdout, "??", path)
	}

	err = idx.WriteUpdates()
	return
}

func showTree(db object.Database, rootOID string, pathPrefix string) (err error) {
	treeObj, err := db.Read(rootOID)
	if err != nil {
		return fmt.Errorf("error reading/parsing head tree: %w", err)
	}

	headTree, err := tree.DeserializeTree(treeObj.Serialize())
	if err != nil {
		return fmt.Errorf("error deserializing head tree: %w", err)
	}

	for _, e := range headTree.Entries() {
		if e.ModeString() == "40000" {
			err = showTree(db, e.OID(), path.Join(pathPrefix, e.Name()))
		} else {
			fmt.Printf("%8s %s %s\n", e.ModeString(), e.OID(), path.Join(pathPrefix, e.Name()))
		}
	}

	return nil
}

// def show_tree(repo, oid, prefix = Pathname.new(""))
// tree = repo.database.load(oid)
// tree.entries.each do |name, entry| path = prefix.join(name)
// if entry.tree?
// show_tree(repo, entry.oid, path) else
// mode = entry.mode.to_s(8)
// puts "#{ mode } #{ entry.oid } #{ path }" end
// end end

func porcelainStatus(bitfield int) (result string) {
	if bitfield&statusindexModified > 0 {
		result = "M"
	} else if bitfield&statusIndexDeleted > 0 {
		result = "D"
	} else {
		result = " "
	}

	if bitfield&statusWorkspaceModified > 0 {
		result += "M"
	} else if bitfield&statusWorkspaceDeleted > 0 {
		result += "D"
	} else {
		result += " "
	}

	return
}
