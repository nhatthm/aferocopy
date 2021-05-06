// +build darwin

package aferocopy

import (
	"os"
	"time"
)

func getTimeSpec(info os.FileInfo) timeSpec {
	stat := fileInfoStat(info.Sys())

	times := timeSpec{
		Mtime: info.ModTime(),
		Atime: time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec),
		Ctime: time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec),
	}

	return times
}
