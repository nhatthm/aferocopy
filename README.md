# `aferocopy`

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/aferocopy)](https://github.com/nhatthm/aferocopy/releases/latest)
[![Build Status](https://github.com/nhatthm/aferocopy/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/aferocopy/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/aferocopy/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/aferocopy)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/aferocopy/v2)](https://goreportcard.com/report/go.nhat.io/aferocopy/v2)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/aferocopy/v2)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

`aferocopy` copies directories recursively using [spf13/afero](https://github.com/spf13/afero)

The idea and logic is ported from [otiai10/copy](https://github.com/otiai10/copy)

## Prerequisites

- `Go >= 1.17`

## Install

```bash
go get go.nhat.io/aferocopy/v2
```

## Usage

```go
package main

import (
	"fmt"

	"go.nhat.io/aferocopy/v2"
)

func main() {
	err := aferocopy.Copy("your/src", "your/dest", aferocopy.Options{
		// Specify the source and destination fs of your choice, default is afero.OsFs.
		// SrcFs: ...,
		// DestFs: ...,
	})

	fmt.Println(err) // nil
}
```
## Advanced Usage

```go
// Options specifies optional actions on copying.
type Options struct {
	// Source filesystem. Default is afero.NewOsFs().
	SrcFs afero.Fs

	// Source filesystem. Default is Options.SrcFs.
	DestFs afero.Fs

	// OnSymlink can specify what to do on symlink
	OnSymlink func(src string) SymlinkAction

	// OnDirExists can specify what to do when there is a directory already existing in destination.
	OnDirExists func(src, dest string) DirExistsAction

	// Skip can specify which files should be skipped
	Skip func(src string) (bool, error)

	// PermissionControl can control permission of
	// every entry.
	// When you want to add permission 0222, do like
	//
	//		PermissionControl = AddPermission(0222)
	//
	// or if you even don't want to touch permission,
	//
	//		PermissionControl = DoNothing
	//
	// By default, PermissionControl = PreservePermission
	PermissionControl PermissionControlFunc

	// Sync file after copy.
	// Useful in case when file must be on the disk
	// (in case crash happens, for example),
	// at the expense of some performance penalty
	Sync bool

	// Preserve the atime and the mtime of the entries
	// On linux we can preserve only up to 1 millisecond accuracy
	PreserveTimes bool

	// Preserve the uid and the gid of all entries.
	PreserveOwner bool

	// The byte size of the buffer to use for copying files.
	// If zero, the internal default buffer of 32KB is used.
	// See https://golang.org/pkg/io/#CopyBuffer for more information.
	CopyBufferSize uint
}
```

```go
package main

import (
	"fmt"
	"strings"

	"go.nhat.io/aferocopy/v2"
)

func main() {
	err := aferocopy.Copy("your/src", "your/dest", aferocopy.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
	})

	fmt.Println(err) // nil
}
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
