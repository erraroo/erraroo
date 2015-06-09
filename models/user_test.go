package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCanSee(t *testing.T) {
	user1 := User{AccountID: 1}
	user2 := User{AccountID: 2}
	project1 := Project{AccountID: 1}
	project2 := Project{AccountID: 2}

	assert.True(t, user1.CanSee(project1))
	assert.False(t, user1.CanSee(project2))

	assert.True(t, user2.CanSee(project2))
	assert.False(t, user2.CanSee(project1))
}
