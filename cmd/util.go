package cmd

import (
	"path"
	"path/filepath"
	"strings"
)

func toAbsolutePath(p string) string {
	if !isAbsolute(p) {
		p = path.Join(wd, p)
	}

	return p
}

func isAbsolute(p string) bool {
	return strings.HasPrefix(p, wd)
}

func toRelativePath(p string) string {
	if isAbsolute(p) {
		p, _ = filepath.Rel(wd, p)
	}

	return p
}
