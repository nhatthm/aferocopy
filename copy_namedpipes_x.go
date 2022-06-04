//go:build windows || plan9 || netbsd || aix || illumos || solaris || js
// +build windows plan9 netbsd aix illumos solaris js

package aferocopy

import (
	"os"

	"github.com/spf13/afero"
)

// copyPipe is for just named pipes. Windows doesn't support them.
func copyPipe(destFs afero.Fs, dest string, info os.FileInfo) error {
	return nil
}
