package aferocopy_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.nhat.io/aferomock"

	"go.nhat.io/aferocopy/v2"
)

func TestPermissionControl_AddPermission_File(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0111
		fileInfo.On("Mode").Return(0o111)
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Expected original + new permissions.
		fs.On("Chmod", "foo.bar", os.FileMode(0o111|0o321)).Return(nil)
	})(t)

	// Set temporary permissions.
	cb, err := aferocopy.AddPermission(0o321)(srcInfo, destFs, "foo.bar")
	require.NoError(t, err)

	// Set final permissions.
	cb(&err)
	require.NoError(t, err)
}

func TestPermissionControl_AddPermission_Dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0111.
		fileInfo.On("Mode").Return(0o111)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Expected original + new permissions
		fs.On("MkdirAll", "foo", os.FileMode(0o755)).Return(nil)
		fs.On("Chmod", "foo", os.FileMode(0o111|0o321)).Return(nil)
	})(t)

	// Set temporary permissions.
	cb, err := aferocopy.AddPermission(0o321)(srcInfo, destFs, "foo")
	require.NoError(t, err)

	// Set final permissions.
	cb(&err)
	require.NoError(t, err)
}

func TestPermissionControl_PreservePermission_File(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123.
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original.
		fs.On("Chmod", "foo.bar", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions.
	cb, err := aferocopy.PreservePermission(srcInfo, destFs, "foo.bar")
	require.NoError(t, err)

	// Set final permissions.
	cb(&err)
	require.NoError(t, err)
}

func TestPermissionControl_PreservePermission_Dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123.
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original.
		fs.On("MkdirAll", "foo", os.FileMode(0o755)).Return(nil)
		fs.On("Chmod", "foo", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions.
	cb, err := aferocopy.PreservePermission(srcInfo, destFs, "foo")
	require.NoError(t, err)

	// Set final permissions.
	cb(&err)
	require.NoError(t, err)
}

func TestPermissionControl_DoNothing_File(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		fileInfo.On("IsDir").Return(false)
	})(t)
	destFs := aferomock.MockFs()(t)

	// Set temporary permissions.
	cb, err := aferocopy.DoNothing(srcInfo, destFs, "foo.bar")
	require.NoError(t, err)

	// Set final permissions.
	cb(&err)
	require.NoError(t, err)
}

func TestPermissionControl_DoNothing_Dir(t *testing.T) {
	t.Parallel()

	// Mocked file info and FS.
	srcInfo := aferomock.MockFileInfo(func(fileInfo *aferomock.FileInfo) {
		// Original permissions 0123.
		fileInfo.On("Mode").Return(0o123)
		fileInfo.On("IsDir").Return(true)
	})(t)
	destFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		// Same permissions as original.
		fs.On("MkdirAll", "foo", os.FileMode(0o123)).Return(nil)
	})(t)

	// Set temporary permissions.
	cb, err := aferocopy.DoNothing(srcInfo, destFs, "foo")
	require.NoError(t, err)

	// Set final permissions
	cb(&err)
	require.NoError(t, err)
}
