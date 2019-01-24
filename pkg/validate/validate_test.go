package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEUCarLicenseNumber(t *testing.T) {
	assert.True(t, EURCarLicenseNumber("AA1234BB"))
	assert.True(t, EURCarLicenseNumber("AI1234OX"))

	assert.False(t, EURCarLicenseNumber("12AI12BP"))
	assert.False(t, EURCarLicenseNumber("AI12BP64"))

	assert.False(t, EURCarLicenseNumber("1234"))
	assert.False(t, EURCarLicenseNumber("1234AB"))
	assert.False(t, EURCarLicenseNumber("AB1234"))
}
