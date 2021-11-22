//go:build !windows && !plan9 && !netbsd && !aix && !illumos && !solaris
// +build !windows,!plan9,!netbsd,!aix,!illumos,!solaris

package aferocopy

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/afero"
)

// copyPipe is for just named pipes.
func copyPipe(destFs afero.Fs, dest string, info os.FileInfo) error {
	if err := destFs.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	return syscall.Mkfifo(dest, uint32(info.Mode()))
}
