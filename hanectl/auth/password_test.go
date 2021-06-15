package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	password := "secret"
	hashed, err := HashPassword(password)
	assert.Nil(t, err)
	assert.True(t, CheckPasswordHash(password, hashed))
}

func TestPasswordHashFail(t *testing.T) {
	password := "secret"
	hashed, err := HashPassword(password)
	assert.Nil(t, err)
	assert.False(t, CheckPasswordHash("no-secret", hashed))
}