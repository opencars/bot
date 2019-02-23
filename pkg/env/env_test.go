package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustFetch(t *testing.T) {
	os.Setenv("VALID_ENV", "test")
	assert.Equal(t, "test", MustFetch("VALID_ENV"))

	assert.Panics(t, func() { MustFetch("INVALID_ENV") })
}

func TestFetch(t *testing.T) {
	t.Run("returns fallback", func(t *testing.T) {
		res := Fetch("INVALID_ENV", "default")
		assert.Equal(t, "default", res)
	})

	t.Run("returns env", func(t *testing.T) {
		os.Setenv("VALID_ENV", "test")
		res := Fetch("VALID_ENV", "default")
		assert.Equal(t, "test", res)
	})
}
