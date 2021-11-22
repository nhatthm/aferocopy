//go:build !windows
// +build !windows

package aferocopy

import (
	"os"
	"syscall"

	"github.com/spf13/afero"
)

func preserveOwner(srcFs afero.Fs, src string, destFs afero.Fs, dest string, info os.FileInfo) (err error) {
	if info == nil {
		if info, err = srcFs.Stat(src); err != nil {
			return err
		}
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if err := destFs.Chown(dest, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}
	}

	return nil
}
