package cmd

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/neocortical/got/repository"
)

func TestInit(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	err := executeInit(initCmd, nil)
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}

	dotGit, err := os.Stat(path.Join(wd, repository.GitDir))
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if !dotGit.IsDir() {
		t.Error("the .git file isn't a directory")
	}

	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}

	expected := regexp.MustCompile(`^Initialized empty Git repository in .*/.git\n$`)
	if !expected.Match(outbuf.Bytes()) {
		t.Errorf("expected output '%s' but got: '%s'", expected.String(), outbuf.String())
	}
}
