package aferocopy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("panic", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			must(errors.New("error"))
		})
	})

	t.Run("not panic", func(t *testing.T) {
		t.Parallel()

		assert.NotPanics(t, func() {
			must(nil)
		})
	})
}
