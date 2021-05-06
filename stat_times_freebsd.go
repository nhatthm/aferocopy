// +build freebsd

package aferocopy

import (
	"os"
	"time"
)

func getTimeSpec(info os.FileInfo) timeSpec {
	stat := fileInfoStat(info.Sys())

	times := timeSpec{
		Mtime: info.ModTime(),
		Atime: time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec)),
		Ctime: time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)),
	}

	return times
}
