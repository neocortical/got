package main

import (
	"fmt"
	"os"

	"github.com/neocortical/got/cmd"
)

const (
	ErrUnknown  = 1
	ErrNotExist = 128
)

func main() {
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)
	cmd.Setenv(os.Getenv)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: unable to get working directory: %v", err)
		os.Exit(1)
	}
	cmd.SetWd(wd)

	cmd.Execute()
}
