package aferocopy

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown(m)

	os.Exit(code)
}

func TestCopy(t *testing.T) {
	err := Copy("./resources/fixtures/data/case00", "./resources/test/data.copy/case00")
	require.NoError(t, err)

	info, err := os.Stat("./resources/test/data.copy/case00/README.md")
	require.NoError(t, err)
	assert.NotNil(t, info)
	assert.False(t, info.IsDir())

	t.Run("specified src does not exist", func(t *testing.T) {
		err := Copy("NOT/EXISTING/SOURCE/PATH", "anywhere")
		require.Error(t, err)
	})

	t.Run("specified src is just a file", func(t *testing.T) {
		err := Copy("resources/fixtures/data/case01/README.md", "resources/test/data.copy/case01/README.md")
		require.NoError(t, err)

		content, err := os.ReadFile("resources/test/data.copy/case01/README.md")

		require.NoError(t, err)
		assert.Equal(t, "case01 - README.md", string(content))
	})

	t.Run("source directory includes symbolic link", func(t *testing.T) {
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03")
		require.NoError(t, err)

		info, err := os.Lstat("resources/test/data.copy/case03/case01")
		require.NoError(t, err)
		assert.NotEqual(t, 0, info.Mode()&os.ModeSymlink)

		t.Run("try to copy to an existing path", func(t *testing.T) {
			err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03")
			require.Error(t, err)
		})
	})

	t.Run("try to copy READ-not-allowed source", func(t *testing.T) {
		err := Copy("resources/fixtures/data/doesNotExist", "resources/test/data.copy/doesNotExist")
		require.Error(t, err)
	})

	t.Run("try to copy a file to existing path", func(t *testing.T) {
		err := Copy("resources/fixtures/data/case04/README.md", "resources/fixtures/data/case04")
		require.Error(t, err)

		err = Copy("resources/fixtures/data/case04/README.md", "resources/fixtures/data/case04/README.md/foobar")
		require.Error(t, err)
	})

	t.Run("try to copy a directory that has no write permission and copy file inside along with it", func(t *testing.T) {
		src := "resources/fixtures/data/case05"
		dest := "resources/test/data.copy/case05"

		err := os.Chmod(src, os.FileMode(0o555))
		require.NoError(t, err)

		err = Copy(src, dest)
		require.NoError(t, err)

		info, err := os.Lstat(dest)
		require.NoError(t, err)

		assert.Equal(t, os.FileMode(0o555), info.Mode().Perm())

		err = os.Chmod(dest, 0o755) //nolint: gosec
		require.NoError(t, err)
	})
}

func TestCopy_NamedPipe(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "js" {
		t.Skip("See https://github.com/otiai10/copy/issues/47")
	}

	t.Run("specified src contains a folder with a named pipe", func(t *testing.T) {
		dest := "resources/test/data.copy/case11"
		err := Copy("resources/fixtures/data/case11", dest)
		require.NoError(t, err)

		info, err := os.Lstat("resources/fixtures/data/case11/foo/bar")
		require.NoError(t, err)
		assert.NotEqual(t, 0, info.Mode()&os.ModeNamedPipe)
		assert.Equal(t, os.FileMode(0o555), info.Mode().Perm())
	})

	t.Run("specified src is a named pipe", func(t *testing.T) {
		dest := "resources/test/data.copy/case11/foo/bar.named"
		err := Copy("resources/fixtures/data/case11/foo/bar", dest)
		require.NoError(t, err)

		info, err := os.Lstat(dest)
		require.NoError(t, err)
		assert.NotEqual(t, 0, info.Mode()&os.ModeNamedPipe)
		assert.Equal(t, os.FileMode(0o555), info.Mode().Perm())
	})
}

func TestOptions_OnSymlink(t *testing.T) {
	t.Run("deep", func(t *testing.T) {
		opt := Options{OnSymlink: func(afero.Fs, string) SymlinkAction { return Deep }}
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03.deep", opt)
		require.NoError(t, err)

		info, err := os.Lstat("resources/test/data.copy/case03.deep/case01")
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0), info.Mode()&os.ModeSymlink)
	})

	t.Run("shallow", func(t *testing.T) {
		opt := Options{OnSymlink: func(afero.Fs, string) SymlinkAction { return Shallow }}
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03.shallow", opt)
		require.NoError(t, err)

		info, err := os.Lstat("resources/test/data.copy/case03.shallow/case01")
		require.NoError(t, err)
		assert.NotEqual(t, os.FileMode(0), info.Mode()&os.ModeSymlink)
	})

	t.Run("skip", func(t *testing.T) {
		opt := Options{OnSymlink: func(afero.Fs, string) SymlinkAction { return Skip }}
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03.skip", opt)
		require.NoError(t, err)

		_, err = os.Stat("resources/test/data.copy/case03.skip/case01")
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("default", func(t *testing.T) {
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03.default")
		require.NoError(t, err)

		info, err := os.Lstat("resources/test/data.copy/case03.default/case01")
		require.NoError(t, err)
		assert.NotEqual(t, os.FileMode(0), info.Mode()&os.ModeSymlink)
	})

	t.Run("not specified", func(t *testing.T) {
		opt := Options{OnSymlink: nil}
		err := Copy("resources/fixtures/data/case03", "resources/test/data.copy/case03.not-specified", opt)
		require.NoError(t, err)

		info, err := os.Lstat("resources/test/data.copy/case03.not-specified/case01")
		require.NoError(t, err)
		assert.NotEqual(t, os.FileMode(0), info.Mode()&os.ModeSymlink)
	})
}

func TestOptions_Skip(t *testing.T) {
	opt := Options{Skip: func(_ afero.Fs, src string) (bool, error) {
		switch {
		case strings.HasSuffix(src, "_skip"):
			return true, nil

		case strings.HasSuffix(src, ".gitfake"):
			return true, nil

		default:
			return false, nil
		}
	}}

	err := Copy("resources/fixtures/data/case06", "resources/test/data.copy/case06", opt)
	require.NoError(t, err)
	info, err := os.Stat("./resources/test/data.copy/case06/dir_skip")
	assert.Nil(t, info)
	assert.True(t, os.IsNotExist(err))

	info, err = os.Stat("./resources/test/data.copy/case06/file_skip")
	assert.Nil(t, info)
	assert.True(t, os.IsNotExist(err))

	info, err = os.Stat("./resources/test/data.copy/case06/README.md")
	assert.NotNil(t, info)
	require.NoError(t, err)

	info, err = os.Stat("./resources/test/data.copy/case06/repo/.gitfake")
	assert.Nil(t, info)
	assert.True(t, os.IsNotExist(err))

	info, err = os.Stat("./resources/test/data.copy/case06/repo/README.md")
	assert.NotNil(t, info)
	require.NoError(t, err)

	t.Run("if Skip func returns error, Copy should be interrupted", func(t *testing.T) {
		errInsideSkipFunc := errors.New("something wrong inside Skip")
		opt := Options{Skip: func(afero.Fs, string) (bool, error) {
			return false, errInsideSkipFunc
		}}
		err := Copy("resources/fixtures/data/case06", "resources/test/data.copy/case06.01", opt)
		assert.Equal(t, errInsideSkipFunc, err)

		files, err := os.ReadDir("./resources/test/data.copy/case06.01")
		require.NoError(t, err)
		assert.Empty(t, files)
	})
}

func TestOptions_PermissionControl(t *testing.T) {
	info, err := os.Stat("resources/fixtures/data/case07/dir_0555")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o555)|os.ModeDir, info.Mode())

	info, err = os.Stat("resources/fixtures/data/case07/file_0444")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o444), info.Mode())

	opt := Options{PermissionControl: AddPermission(0o222)}
	err = Copy("resources/fixtures/data/case07", "resources/test/data.copy/case07", opt)
	require.NoError(t, err)

	info, err = os.Stat("resources/test/data.copy/case07/dir_0555")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o555|0o222)|os.ModeDir, info.Mode())

	info, err = os.Stat("resources/test/data.copy/case07/file_0444")
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o444|0o222), info.Mode())
}

func TestOptions_Sync(t *testing.T) {
	// With Sync option, each file will be flushed to storage on copying.
	// TODO: Since it's a bit hard to simulate real usecases here. This testcase is nonsense.
	//nolint: godox
	opt := Options{Sync: true}
	err := Copy("resources/fixtures/data/case08", "resources/test/data.copy/case08", opt)
	require.NoError(t, err)
}

func TestOptions_PreserveTimes(t *testing.T) {
	err := Copy("resources/fixtures/data/case09", "resources/test/data.copy/case09")
	require.NoError(t, err)

	opt := Options{PreserveTimes: true}
	err = Copy("resources/fixtures/data/case09", "resources/test/data.copy/case09-preservetimes", opt)
	require.NoError(t, err)

	for _, entry := range []string{"", "README.md", "symlink"} {
		orig, err := os.Stat("resources/fixtures/data/case09/" + entry)
		require.NoError(t, err)

		plain, err := os.Stat("resources/test/data.copy/case09/" + entry)
		require.NoError(t, err)

		preserved, err := os.Stat("resources/test/data.copy/case09-preservetimes/" + entry)
		require.NoError(t, err)

		assert.NotEqual(t, orig.ModTime().Unix(), plain.ModTime().Unix())
		assert.Equal(t, orig.ModTime().Unix(), preserved.ModTime().Unix())
	}
}

func TestOptions_OnDirExists(t *testing.T) {
	err := Copy("resources/fixtures/data/case10/dest", "resources/test/data.copy/case10/dest.1")
	require.NoError(t, err)

	err = Copy("resources/fixtures/data/case10/dest", "resources/test/data.copy/case10/dest.2")
	require.NoError(t, err)

	err = Copy("resources/fixtures/data/case10/dest", "resources/test/data.copy/case10/dest.3")
	require.NoError(t, err)

	t.Run("replace", func(t *testing.T) {
		opt := Options{
			OnDirExists: func(afero.Fs, string, afero.Fs, string) DirExistsAction {
				return Merge
			},
		}

		err := Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.1", opt)
		require.NoError(t, err)

		err = Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.1", opt)
		require.NoError(t, err)

		b, err := os.ReadFile("resources/test/data.copy/case10/dest.1/" + "foo/" + "text_aaa")
		require.NoError(t, err)
		assert.Equal(t, "This is text_aaa from src", string(b))

		stat, err := os.Stat("resources/test/data.copy/case10/dest.1/foo/text_eee")
		require.NoError(t, err)
		assert.NotNil(t, stat)
	})

	t.Run("replace", func(t *testing.T) {
		opt := Options{
			OnDirExists: func(afero.Fs, string, afero.Fs, string) DirExistsAction {
				return Replace
			},
		}
		err := Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.2", opt)
		require.NoError(t, err)

		err = Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.2", opt)
		require.NoError(t, err)

		b, err := os.ReadFile("resources/test/data.copy/case10/dest.2/" + "foo/" + "text_aaa")
		require.NoError(t, err)
		assert.Equal(t, "This is text_aaa from src", string(b))

		stat, err := os.Stat("resources/test/data.copy/case10/dest.2/foo/text_eee")
		assert.True(t, os.IsNotExist(err))
		assert.Nil(t, stat)
	})

	t.Run("untouchable", func(t *testing.T) {
		opt := Options{
			OnDirExists: func(afero.Fs, string, afero.Fs, string) DirExistsAction {
				return Untouchable
			},
		}
		err := Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.3", opt)
		require.NoError(t, err)

		b, err := os.ReadFile("resources/test/data.copy/case10/dest.3/" + "foo/" + "text_aaa")
		require.NoError(t, err)
		assert.Equal(t, "This is text_aaa from dest", string(b))
	})

	t.Run("PreserveTimes is true with Untouchable", func(t *testing.T) {
		opt := Options{
			OnDirExists: func(afero.Fs, string, afero.Fs, string) DirExistsAction {
				return Untouchable
			},
			PreserveTimes: true,
		}
		err := Copy("resources/fixtures/data/case10/src", "resources/test/data.copy/case10/dest.3", opt)
		require.NoError(t, err)
	})
}

func TestOptions_CopyBufferSize(t *testing.T) {
	opt := Options{
		CopyBufferSize: 512,
	}

	err := Copy("resources/fixtures/data/case12", "resources/test/data.copy/case12", opt)

	require.NoError(t, err)

	content, err := os.ReadFile("resources/test/data.copy/case12/README.md")

	require.NoError(t, err)
	assert.Equal(t, "case12 - README.md", string(content))
}

func TestOptions_PreserveOwner(t *testing.T) {
	opt := Options{
		PreserveOwner: true,
	}

	err := Copy("resources/fixtures/data/case13", "resources/test/data.copy/case13", opt)
	require.NoError(t, err)
}
