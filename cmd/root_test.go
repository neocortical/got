package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func setUpTestWorkspace(t *testing.T, env map[string]string) (outbuf, errbuf *bytes.Buffer) {
	outbuf = new(bytes.Buffer)
	errbuf = new(bytes.Buffer)
	stdout = outbuf
	stderr = errbuf

	tempdir, err := ioutil.TempDir("", "got_test_*")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	wd = tempdir

	if env == nil {
		env = map[string]string{}
	}
	getenv = func(key string) string {
		return env[key]
	}

	return
}

func tearDownTestWorkspace() {
	stdout = nil
	stderr = nil
	getenv = nil
	err := os.RemoveAll(wd)
	if err != nil {
		fmt.Println("error deleting temp workspace:", err)
	}
	wd = ""
}

func writeFile(t *testing.T, filename string) {
	dir, _ := filepath.Split(filename)
	if dir != "" {
		mkdir(t, dir)
	}

	err := ioutil.WriteFile(path.Join(wd, filename), []byte{}, 0644)
	if err != nil {
		t.Fatalf("error writing file: %v", err)
	}
}

func mkdir(t *testing.T, dirname string) {
	err := os.MkdirAll(path.Join(wd, dirname), 0755)
	if err != nil {
		t.Fatalf("error creating dir '%s': %v", dirname, err)
	}
}

func initOrDie(t *testing.T) {
	err := executeInit(initCmd, nil)
	if err != nil {
		t.Fatalf("expected no errors during init but got: %v", err)
	}
}
