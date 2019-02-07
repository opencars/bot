package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEuroPlates(t *testing.T) {
	assert.True(t, EuroPlates("AA1111XX"))

	assert.False(t, EuroPlates("AA111XX1"))
	assert.False(t, EuroPlates("1AA111XX"))
	assert.False(t, EuroPlates("1AA111XXX"))
	assert.False(t, EuroPlates("1AA111XX1"))
	assert.False(t, EuroPlates("Test with number AA1111XX inside"))
}
