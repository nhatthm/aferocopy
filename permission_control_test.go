package aferocopy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.nhat.io/aferomock"
)

func TestPermissionControl_AddPermission_file(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0111
		fileInfo.On("Mode").Return(0o111)
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Expected original + new permissions
		fs.On("Chmod", "foo.bar", os.FileMode(0o111|0o321)).Return(nil)
	})(t)

	// Set temporary permissions
	cb, err := AddPermission(0o321)(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}

func TestPermissionControl_AddPermission_dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0111
		fileInfo.On("Mode").Return(0o111)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Expected original + new permissions
		fs.On("MkdirAll", "foo.bar", tmpPermissionForDirectory).Return(nil)
		fs.On("Chmod", "foo.bar", os.FileMode(0o111|0o321)).Return(nil)
	})(t)

	// Set temporary permissions
	cb, err := AddPermission(0o321)(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}

func TestPermissionControl_PreservePermission_file(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original
		fs.On("Chmod", "foo.bar", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions
	cb, err := PreservePermission(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}

func TestPermissionControl_PreservePermission_dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original
		fs.On("MkdirAll", "foo.bar", tmpPermissionForDirectory).Return(nil)
		fs.On("Chmod", "foo.bar", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions
	cb, err := PreservePermission(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}

func TestPermissionControl_DoNothing_file(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs()(t)

	// Set temporary permissions
	cb, err := DoNothing(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}

func TestPermissionControl_DoNothing_dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original
		fs.On("MkdirAll", "foo.bar", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions
	cb, err := DoNothing(srcInfo, destFs, "foo.bar")
	assert.NoError(t, err)

	// Set final permissions
	cb(&err)
	assert.NoError(t, err)
}
