package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateExceeded(t *testing.T) {
	_, _, project := aup(t)
	err := RateExceeded(project.Token)
	assert.Nil(t, err)
}
