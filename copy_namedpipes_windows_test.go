//go:build windows
// +build windows

package aferocopy

import (
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/stretchr/testify/assert"
)

func TestCopyPipe_CouldNotMkdir(t *testing.T) {
	t.Parallel()

	fs := aferomock.NoMockFs(t)

	err := copyPipe(fs, "/path/to/pipe", aferomock.NewFileInfo())
	assert.NoError(t, err)
}
