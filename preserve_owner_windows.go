//go:build windows
// +build windows

package aferocopy

import (
	"os"

	"github.com/spf13/afero"
)

func preserveOwner(srcFs afero.Fs, src string, destFs afero.Fs, dest string, info os.FileInfo) (err error) {
	return nil
}
