package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "got",
		Short: "A VCS.",
		Long:  `got is a clone of git, which is a little-known version control system.`,
	}

	stdout io.Writer
	stderr io.Writer
	getenv func(string) string
	wd     string
)

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(statusCmd)
}

func SetStdout(w io.Writer) {
	stdout = w
}

func SetStderr(w io.Writer) {
	stderr = w
}

func Setenv(f func(string) string) {
	getenv = f
}

func SetWd(dir string) {
	wd = dir
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(stderr, fmt.Sprintf("Fatal: %v", err))
		os.Exit(1)
	}
}
