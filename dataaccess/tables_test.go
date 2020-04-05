package dataaccess

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAppUser(t *testing.T) {
	u := &AppUser{Username: "admin", Password: "admin"}
	err := u.Validate(true, false, false)
	assert.Nil(t, err)

	// 33文字以上
	u.Username = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	err = u.Validate(true, true, false)
	assert.EqualError(t, err, "username must be less than or equals to 32 characters")

	err = u.Validate(false, true, false)
	assert.Nil(t, err)

	// 英数字以外
	u.Username = "abcd1234_"
	err = u.Validate(true, true, false)
	assert.EqualError(t, err, "username must contain alphabet or numeric")

	// 空文字
	u.Username = ""
	err = u.Validate(true, false, false)
	assert.EqualError(t, err, "username must not be empty")

	// 6文字未満
	u.Username = "abcde"
	err = u.Validate(true, true, false)
	assert.EqualError(t, err, "username must be greater than 5 characters")

	// 33文字以上
	u.Password = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	err = u.Validate(false, true, true)
	assert.EqualError(t, err, "password must be less than or equals to 32 characters")

	// 空文字
	u.Password = ""
	err = u.Validate(false, true, false)
	assert.EqualError(t, err, "password must not be empty")

	// 8文字未満
	u.Password = "abcd123"
	err = u.Validate(false, true, true)
	assert.EqualError(t, err, "password must be greater than 7 characters")
}
