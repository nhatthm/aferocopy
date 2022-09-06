package aferocopy

import (
	"os"

	"github.com/spf13/afero"
)

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0o755)
)

// PermissionControlFunc is a function that can be used to control the permission of a file or directory while copying.
type PermissionControlFunc func(srcInfo os.FileInfo, destFs afero.Fs, dest string) (chmodFunc func(*error), err error)

var (
	// AddPermission controls the permission of the destination file.
	AddPermission = func(perm os.FileMode) PermissionControlFunc {
		return func(srcInfo os.FileInfo, destFs afero.Fs, dest string) (func(*error), error) {
			orig := srcInfo.Mode()

			if srcInfo.IsDir() {
				if err := destFs.MkdirAll(dest, tmpPermissionForDirectory); err != nil {
					return func(*error) {}, err
				}
			}

			return func(err *error) {
				chmod(destFs, dest, orig|perm, err)
			}, nil
		}
	}

	// PreservePermission preserves the original permission.
	PreservePermission = AddPermission(0)

	// DoNothing do not touch the permission.
	DoNothing = PermissionControlFunc(func(srcInfo os.FileInfo, destFs afero.Fs, dest string) (func(*error), error) {
		if srcInfo.IsDir() {
			if err := destFs.MkdirAll(dest, srcInfo.Mode()); err != nil {
				return func(*error) {}, err
			}
		}

		return func(*error) {}, nil
	})
)

// chmod ANYHOW changes file mode,
// with assigning error raised during Chmod,
// BUT respecting the error already reported.
func chmod(fs afero.Fs, dir string, mode os.FileMode, reported *error) {
	if err := fs.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}
