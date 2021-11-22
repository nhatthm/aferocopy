//go:build plan9 || netbsd
// +build plan9 netbsd

package aferocopy

import (
	"os"
)

func getTimeSpec(info os.FileInfo) timespec {
	times := timespec{
		Mtime: info.ModTime(),
		Atime: info.ModTime(),
		Ctime: info.ModTime(),
	}

	return times
}
