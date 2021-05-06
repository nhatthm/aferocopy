// +build !go1.16

package aferocopy

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy_PathError(t *testing.T) {
	t.Run("too long name is given", func(t *testing.T) {
		dest := "foobar"

		for i := 0; i < 8; i++ {
			dest = dest + dest
		}

		err := Copy("resources/fixtures/data/case00", filepath.Join("resources/test/data/case00", dest))

		assert.NotNil(t, err)
		assert.IsType(t, &os.PathError{}, err)
	})

	t.Run("try to create not permitted location", func(t *testing.T) {
		if runtime.GOOS == "windows" || runtime.GOOS == "freebsd" || os.Getenv("TESTCASE") != "" {
			t.Skipf("FIXME: error IS nil here in Windows and FreeBSD")
		}

		err := Copy("resources/fixtures/data/case00", "/case00")

		assert.NotNil(t, err)
		assert.IsType(t, &os.PathError{}, err)
	})

	t.Run("try to create a directory on existing file name", func(t *testing.T) {
		err := Copy("resources/fixtures/data/case02", "resources/test/data.copy/case00/README.md")

		assert.NotNil(t, err)
		assert.IsType(t, &os.PathError{}, err)
	})
}
