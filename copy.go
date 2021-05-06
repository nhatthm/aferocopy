package aferocopy

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0755)
)

type timeSpec struct {
	Mtime time.Time
	Atime time.Time
	Ctime time.Time
}

// Copy copies src to dest, doesn't matter if src is a directory or a file.
func Copy(src, dest string, opt ...Options) error {
	o := assure(src, dest, opt...)

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
func switchboard(src, dest string, info os.FileInfo, opt Options) (err error) {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		err = onSymlink(src, dest, opt)

	case info.IsDir():
		err = copyDir(src, dest, info, opt)

	case info.Mode()&os.ModeNamedPipe != 0:
		err = copyPipe(opt.DestFs, dest, info)

	default:
		err = copyFile(src, dest, info, opt)
	}

	return err
}

// copy decide if this src should be copied or not.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
// nolint: predeclared
//goland:noinspection GoReservedWordUsedAsName
func copy(src, dest string, info os.FileInfo, opt Options) error {
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

	if err = destFs.Chmod(f.Name(), info.Mode()|opt.AddPermission); err != nil {
		return
	}

	s, err := srcFs.Open(src)
	if err != nil {
		return
	}

	defer closeFile(s, &err)

	if _, err = io.Copy(f, s); err != nil {
		return
	}

	if opt.Sync {
		err = f.Sync()
	}

	if opt.PreserveTimes {
		return preserveTimes(info, destFs, dest)
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
func copyDir(srcDir, destDir string, info os.FileInfo, opt Options) (err error) {
	srcFs := opt.SrcFs
	destFs := opt.DestFs

	exit, err := checkDir(srcDir, destDir, opt)
	if err != nil || exit {
		return err
	}

	originalMode := info.Mode()

	// Make dest dir with 0755 so that everything writable.
	if err = destFs.MkdirAll(destDir, tmpPermissionForDirectory); err != nil {
		return
	}
	// Recover dir mode with original one.
	defer chmod(destFs, destDir, originalMode|opt.AddPermission, &err)

	contents, err := afero.ReadDir(srcFs, srcDir)
	if err != nil {
		return
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcDir, content.Name()), filepath.Join(destDir, content.Name())

		if err = copy(cs, cd, content, opt); err != nil {
			// If any error, exit immediately.
			return
		}
	}

	if opt.PreserveTimes {
		return preserveTimes(info, destFs, destDir)
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

		return copy(orig, dest, info, opt)

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

// chmod ANYHOW changes file mode,
// with assigning error raised during Chmod,
// BUT respecting the error already reported.
func chmod(fs afero.Fs, dir string, mode os.FileMode, reported *error) {
	if err := fs.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}

// assure Options struct, should be called only once.
// All optional values MUST NOT BE nil/zero after assured.
func assure(src, dest string, opts ...Options) Options {
	defaults := getDefaultOptions(src, dest)
	defaults.SrcFs = afero.NewOsFs()
	defaults.DestFs = defaults.SrcFs

	if len(opts) == 0 {
		return defaults
	}

	if opts[0].SrcFs == nil {
		opts[0].SrcFs = defaults.SrcFs
	}

	if opts[0].DestFs == nil {
		opts[0].DestFs = opts[0].SrcFs
	}

	if opts[0].OnSymlink == nil {
		opts[0].OnSymlink = defaults.OnSymlink
	}

	if opts[0].Skip == nil {
		opts[0].Skip = defaults.Skip
	}

	opts[0].intent.src = defaults.intent.src
	opts[0].intent.dest = defaults.intent.dest

	return opts[0]
}
