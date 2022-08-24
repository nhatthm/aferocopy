//go:build !windows
// +build !windows

package aferocopy

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.nhat.io/aferomock"
)

func TestCopyPipe_CouldNotMkdir(t *testing.T) {
	t.Parallel()

	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("MkdirAll", "/path/to", os.ModePerm).
			Return(errors.New("could not mkdir"))
	})(t)

	err := copyPipe(fs, "/path/to/pipe", aferomock.NewFileInfo())

	expectedErr := `could not mkdir`

	assert.EqualError(t, err, expectedErr)
}
