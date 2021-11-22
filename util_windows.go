//go:build windows
// +build windows

package aferocopy

import "syscall"

func fileInfoStat(v interface{}) *syscall.Win32FileAttributeData {
	s, ok := v.(*syscall.Win32FileAttributeData)
	if !ok {
		panic("not a *syscall.Win32FileAttributeData")
	}

	return s
}
