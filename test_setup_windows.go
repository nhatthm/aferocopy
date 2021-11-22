//go:build windows
// +build windows

package aferocopy

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func setup(m *testing.M) {
	fs := afero.NewOsFs()

	must(fs.MkdirAll("resources/test/data.copy", os.ModePerm))
	must(fs.Chmod("resources/fixtures/data/case07/dir_0555", 0555))
	must(fs.Chmod("resources/fixtures/data/case07/file_0444", 0444))

	if fs, ok := fs.(afero.Linker); ok {
		must(fs.SymlinkIfPossible("resources/fixtures/data/case01", "resources/fixtures/data/case03/case01"))
	}
}

func teardown(*testing.M) {
	fs := afero.NewOsFs()

	must(fs.RemoveAll("resources/test"))
	must(fs.Remove("resources/fixtures/data/case03/case01"))
}
