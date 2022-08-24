//go:build windows || plan9 || netbsd || aix || illumos || solaris || js
// +build windows plan9 netbsd aix illumos solaris js

package aferocopy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.nhat.io/aferomock"
)

func TestCopyPipe_CouldNotMkdir(t *testing.T) {
	t.Parallel()

	fs := aferomock.NoMockFs(t)

	err := copyPipe(fs, "/path/to/pipe", aferomock.NewFileInfo())
	assert.NoError(t, err)
}
