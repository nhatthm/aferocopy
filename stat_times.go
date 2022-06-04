//go:build !windows && !darwin && !freebsd && !plan9 && !netbsd && !js
// +build !windows,!darwin,!freebsd,!plan9,!netbsd,!js

package aferocopy

import (
	"os"
	"time"
)

func getTimeSpec(info os.FileInfo) timeSpec {
	stat := fileInfoStat(info.Sys())

	times := timeSpec{
		Mtime: info.ModTime(),
		Atime: time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)),
		Ctime: time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
	}

	return times
}
