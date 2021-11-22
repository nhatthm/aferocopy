//go:build windows
// +build windows

package aferocopy

import (
	"os"
	"time"
)

func getTimeSpec(info os.FileInfo) timeSpec {
	stat := fileInfoStat(info.Sys())

	return timeSpec{
		Mtime: time.Unix(0, stat.LastWriteTime.Nanoseconds()),
		Atime: time.Unix(0, stat.LastAccessTime.Nanoseconds()),
		Ctime: time.Unix(0, stat.CreationTime.Nanoseconds()),
	}
}
