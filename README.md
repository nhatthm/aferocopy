# `aferocopy`

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/aferocopy)](https://github.com/nhatthm/aferocopy/releases/latest)
[![Build Status](https://github.com/nhatthm/aferocopy/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/aferocopy/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/aferocopy/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/aferocopy)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/aferocopy)](https://goreportcard.com/report/github.com/nhatthm/aferocopy)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/aferocopy)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

`aferocopy` copies directories recursively using [spf13/afero](https://github.com/spf13/afero)

The idea and logic is ported from [otiai10/copy](https://github.com/otiai10/copy)

## Prerequisites

- `Go >= 1.16`

## Install

```bash
go get github.com/nhatthm/aferocopy
```

## Usage

```go
err := aferocopy.Copy("your/directory", "your/directory.copy")
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

	// AddPermission to every entry,
	// NO MORE THAN 0777
	AddPermission os.FileMode

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
// For example...
opt := Options{
	Skip: func(src string) (bool, error) {
		return strings.HasSuffix(src, ".git"), nil
	},
}
err := Copy("your/directory", "your/directory.copy", opt)
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
