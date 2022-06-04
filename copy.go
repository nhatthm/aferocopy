package aferocopy

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

type timeSpec struct {
	Mtime time.Time
	Atime time.Time
	Ctime time.Time
}

// Copy copies src to dest, doesn't matter if src is a directory or a file.
func Copy(src, dest string, opt ...Options) error {
	o := assureOptions(src, dest, opt...)

	info, err := stat(o.SrcFs, src)
	if err != nil {
		return err
	}

	return switchboard(src, dest, info, o)
}

func stat(fs afero.Fs, path string) (os.FileInfo, error) {
	if fs, ok := fs.(afero.Lstater); ok {
		fi, _, err := fs.LstatIfPossible(path)

		return fi, err
	}

	return fs.Stat(path)
}

// switchboard switches proper copy functions regarding file type, etc...
// If there would be anything else here, add a case to this switchboard.
func switchboard(src, dest string, info os.FileInfo, opt Options) error {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return onSymlink(src, dest, opt)

	case info.IsDir():
		return copyDir(src, dest, info, opt)

	case info.Mode()&os.ModeNamedPipe != 0:
		return copyPipe(opt.DestFs, dest, info)

	default:
		return copyFile(src, dest, info, opt)
	}
}

// copyNextOrSkip decide if this src should be copied or not.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copyNextOrSkip(src, dest string, info os.FileInfo, opt Options) error {
	skip, err := opt.Skip(opt.SrcFs, src)
	if err != nil {
		return err
	}

	if skip {
		return nil
	}

	return switchboard(src, dest, info, opt)
}

// copyFile is for just a file,
// with considering existence of parent directory
// and file permission.
// nolint: cyclop
func copyFile(src, dest string, info os.FileInfo, opt Options) (err error) {
	srcFs := opt.SrcFs
	destFs := opt.DestFs

	if err = destFs.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return
	}

	f, err := destFs.Create(dest)
	if err != nil {
		return
	}

	defer closeFile(f, &err)

	chmod, err := opt.PermissionControl(info, destFs, dest)
	if err != nil {
		return err
	}

	chmod(&err)

	s, err := srcFs.Open(src)
	if err != nil {
		return
	}

	defer closeFile(s, &err)

	var (
		buf []byte
		w   io.Writer = f
	)

	if opt.CopyBufferSize != 0 {
		buf = make([]byte, opt.CopyBufferSize)
		// Disable using `ReadFrom` by io.CopyBuffer.
		w = struct{ io.Writer }{f}
	}

	if _, err = io.CopyBuffer(w, s, buf); err != nil {
		return err
	}

	if opt.Sync {
		err = f.Sync()
	}

	if opt.PreserveOwner {
		if err := preserveOwner(srcFs, src, destFs, dest, info); err != nil {
			return err
		}
	}

	if opt.PreserveTimes {
		if err := preserveTimes(info, destFs, dest); err != nil {
			return err
		}
	}

	return nil
}

func checkDir(srcDir, destDir string, opt Options) (exit bool, err error) {
	srcFs := opt.SrcFs
	destFs := opt.DestFs

	_, err = destFs.Stat(destDir)
	if err == nil && opt.OnDirExists != nil && destDir != opt.intent.dest {
		switch opt.OnDirExists(srcFs, srcDir, destFs, destDir) {
		case Replace:
			if err := destFs.RemoveAll(destDir); err != nil {
				return true, err
			}

		case Untouchable:
			return true, nil

		// case "Merge" is default behavior. Go through.
		case Merge:
			return false, nil
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return true, err // Unwelcome error type...!
	}

	return false, nil
}

// copyDir is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
// nolint: cyclop
func copyDir(srcDir, destDir string, info os.FileInfo, opt Options) (err error) {
	srcFs := opt.SrcFs
	destFs := opt.DestFs

	exit, err := checkDir(srcDir, destDir, opt)
	if err != nil || exit {
		return err
	}

	chmod, err := opt.PermissionControl(info, destFs, destDir)
	if err != nil {
		return err
	}

	defer chmod(&err)

	contents, err := afero.ReadDir(srcFs, srcDir)
	if err != nil {
		return
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcDir, content.Name()), filepath.Join(destDir, content.Name())

		if err = copyNextOrSkip(cs, cd, content, opt); err != nil {
			// If any error, exit immediately.
			return
		}
	}

	if opt.PreserveOwner {
		if err := preserveOwner(srcFs, srcDir, destFs, destDir, info); err != nil {
			return err
		}
	}

	if opt.PreserveTimes {
		if err := preserveTimes(info, destFs, destDir); err != nil {
			return err
		}
	}

	return nil
}

func onSymlink(src, dest string, opt Options) error {
	destFs, ok := opt.DestFs.(afero.Symlinker)
	if !ok {
		return afero.ErrNoSymlink
	}

	switch opt.OnSymlink(opt.SrcFs, src) {
	case Shallow:
		return copySymlink(src, destFs, dest)

	case Deep:
		orig, err := destFs.ReadlinkIfPossible(src)
		if err != nil {
			return err
		}

		info, _, err := destFs.LstatIfPossible(orig)
		if err != nil {
			return err
		}

		return copyNextOrSkip(orig, dest, info, opt)

	case Skip:
		fallthrough

	default:
		return nil // do nothing
	}
}

// copySymlink is for a symlink,
// with just creating a new symlink by replicating src symlink.
func copySymlink(src string, destFs afero.Symlinker, dest string) error {
	src, err := destFs.ReadlinkIfPossible(src)
	if err != nil {
		return err
	}

	return destFs.SymlinkIfPossible(src, dest)
}

// closeFile ANYHOW closes file,
// with assigning error raised during Close,
// BUT respecting the error already reported.
func closeFile(f afero.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}
