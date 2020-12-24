package cmd

import (
	"testing"
)

func TestListUntrackedFilesInOrder(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	initOrDie(t)
	outbuf.Reset()
	writeFile(t, "file.txt")
	writeFile(t, "anotherfile.txt")

	err := executeStatus(statusCmd, []string{})
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}
	expected := "?? anotherfile.txt\n?? file.txt\n"

	if outbuf.String() != expected {
		t.Errorf("expected output '%s' but got: '%s'", expected, outbuf.String())
	}
}

func TestListOnlyUntrackedFilesInOrder(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	initOrDie(t)
	writeFile(t, "file1.txt")

	err := executeAdd(addCmd, []string{"file1.txt"})
	if err != nil {
		t.Fatalf("expected no errors during add but got: %v", err)
	}

	writeFile(t, "file2.txt")

	outbuf.Reset()

	err = executeStatus(statusCmd, []string{})
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}
	expected := "?? file2.txt\n"

	if outbuf.String() != expected {
		t.Errorf("expected output '%s' but got: '%s'", expected, outbuf.String())
	}
}

func TestListUntrackedDirectoriesNotContents(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	initOrDie(t)
	outbuf.Reset()
	writeFile(t, "file1.txt")
	writeFile(t, "dir/file2.txt")

	err := executeStatus(statusCmd, []string{})
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}
	expected := "?? dir/\n?? file1.txt\n"

	if outbuf.String() != expected {
		t.Errorf("expected output '%s' but got: '%s'", expected, outbuf.String())
	}
}

func TestListUntrackedFilesInsideTrackedDirectories(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	initOrDie(t)
	outbuf.Reset()
	writeFile(t, "a/b/inner.txt")

	err := executeAdd(addCmd, []string{"."})
	if err != nil {
		t.Fatalf("expected no errors during add but got: %v", err)
	}

	writeFile(t, "a/outer.txt")
	writeFile(t, "a/b/c/file.txt")

	err = executeStatus(statusCmd, []string{})
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}
	expected := "?? a/b/c/\n?? a/outer.txt\n"

	if outbuf.String() != expected {
		t.Errorf("expected output '%s' but got: '%s'", expected, outbuf.String())
	}
}

func TestDontListUntrackedEmptyDirectories(t *testing.T) {
	outbuf, errbuf := setUpTestWorkspace(t, nil)
	defer tearDownTestWorkspace()

	initOrDie(t)
	outbuf.Reset()
	mkdir(t, "untracked")

	err := executeStatus(statusCmd, []string{})
	if err != nil {
		t.Errorf("expected no errors but got: %v", err)
	}
	if errbuf.Len() > 0 {
		t.Errorf("expected no error output but got: %s", errbuf.String())
	}
	expected := ""

	if outbuf.String() != expected {
		t.Errorf("expected output '%s' but got: '%s'", expected, outbuf.String())
	}
}
