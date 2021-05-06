// +build !windows

package aferocopy

import "syscall"

func fileInfoStat(v interface{}) *syscall.Stat_t {
	s, ok := v.(*syscall.Stat_t)
	if !ok {
		panic("not a *syscall.Stat_t")
	}

	return s
}
