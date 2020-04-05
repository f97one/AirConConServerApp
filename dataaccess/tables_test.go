package dataaccess

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAppUser(t *testing.T) {
	u := &AppUser{Username: "admin", Password: "admin"}
	err := u.Validate(true)
	assert.Nil(t, err)

	u.Username = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	err = u.Validate(true)
	assert.EqualError(t, err, "username must be less than or equals to 32 characters")

	err = u.Validate(false)
	assert.Nil(t, err)

	u.Username = "abcd1234_"
	err = u.Validate(true)
	assert.EqualError(t, err, "username must contain alphabet or numeric")

	u.Username = ""
	err = u.Validate(true)
	assert.EqualError(t, err, "username must not be empty")

	u.Password = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	err = u.Validate(false)
	assert.EqualError(t, err, "password must be less than or equals to 32 characters")

	u.Password = ""
	err = u.Validate(false)
	assert.EqualError(t, err, "password must not be empty")
}
