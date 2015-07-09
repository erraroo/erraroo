package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_ShouldNotify(t *testing.T) {
	group := &Error{}

	ok := group.ShouldNotify()
	assert.True(t, ok)

	group.Occurrences = 0
	ok = group.ShouldNotify()
	assert.True(t, ok)

	group.Occurrences = 1
	ok = group.ShouldNotify()
	assert.False(t, ok)

	group.Occurrences = 1
	group.Resolved = true
	ok = group.ShouldNotify()
	assert.True(t, ok)

	group.Muted = true
	ok = group.ShouldNotify()
	assert.False(t, ok)
}
