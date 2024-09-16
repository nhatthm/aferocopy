package aferocopy

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
)

func ExampleCopy() {
	err := Copy("resources/fixtures/data/example", "resources/test/data.copy/example")
	fmt.Println("Error:", err)

	info, _ := os.Stat("resources/test/data.copy/example") //nolint: errcheck
	fmt.Println("IsDir:", info.IsDir())

	// Output:
	// Error: <nil>
	// IsDir: true
}

func ExampleOptions() {
	err := Copy(
		"resources/fixtures/data/example",
		"resources/test/data.copy/example_with_options",
		Options{
			Skip: func(_ afero.Fs, src string) (bool, error) {
				return strings.HasSuffix(src, ".git-like"), nil
			},
			OnSymlink: func(afero.Fs, string) SymlinkAction {
				return Skip
			},
			PermissionControl: AddPermission(0o200),
		},
	)
	fmt.Println("Error:", err)

	_, err = os.Stat("resources/test/data.copy/example_with_options/.git-like")
	fmt.Println("Skipped:", os.IsNotExist(err))

	// Output:
	// Error: <nil>
	// Skipped: true
}
