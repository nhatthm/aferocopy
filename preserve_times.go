package aferocopy

import (
	"os"

	"github.com/spf13/afero"
)

func preserveTimes(srcInfo os.FileInfo, destFs afero.Fs, dest string) error {
	spec := getTimeSpec(srcInfo)

	return destFs.Chtimes(dest, spec.Atime, spec.Mtime)
}
