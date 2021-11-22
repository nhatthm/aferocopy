//go:build !windows
// +build !windows

package aferocopy

import (
	"errors"
	"syscall"
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPreserveOwner_statFail(t *testing.T) {
	t.Parallel()

	srcFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", mock.Anything).
			Return(nil, errors.New("stat error"))
	})(t)

	actual := preserveOwner(srcFs, "resources", nil, "", nil)
	expected := errors.New("stat error")

	assert.Equal(t, expected, actual)
}

func TestPreserveOwner_chownSuccess(t *testing.T) {
	t.Parallel()

	const src = "resources/fixtures/data/case00/README.md"

	srcFs := afero.NewOsFs()
	info, err := srcFs.Stat(src)

	require.NoError(t, err)

	stat, ok := info.Sys().(*syscall.Stat_t)

	require.True(t, ok)

	expectedUID := int(stat.Uid)
	expectedGID := int(stat.Gid)

	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Chown", src, expectedUID, expectedGID).
			Return(nil)
	})(t)

	err = preserveOwner(srcFs, src, destFs, src, nil)

	require.NoError(t, err)
}

func TestPreserveOwner_chownFail(t *testing.T) {
	t.Parallel()

	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Chown", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("chown error"))
	})(t)

	actual := preserveOwner(afero.NewOsFs(), "resources/fixtures/data/case00/README.md", destFs, "", nil)
	expected := errors.New("chown error")

	assert.Equal(t, expected, actual)
}
