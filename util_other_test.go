// +build !windows

package aferocopy

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileInfoStat(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		var input interface{} = &syscall.Stat_t{}

		actual := fileInfoStat(input)
		expected := &syscall.Stat_t{}

		assert.Equal(t, expected, actual)
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Panics(t, func() {
			fileInfoStat(42)
		})
	})
}
