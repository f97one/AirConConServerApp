package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateScheduleResp(t *testing.T) {
	r := &scheduleResp{
		ScheduleId: "",
		Name:       "",
		OnOff:      "",
		Weekday:    []int{1},
		Time:       "00:00",
		ScriptId:   "d1a0e190acc9f55ef7032dedd1a0eeaf98591ac7",
	}

	// name
	// 33文字以上
	r.Name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"
	err := r.validate()
	assert.EqualError(t, err, "name must be less than or equals to 32 characters")

	// 文字種
	r.Name = "1234abcdABCD_/"
	err = r.validate()
	assert.EqualError(t, err, "name must contain alphabet or number or underscore only")

	// onOff
	r.Name = "abcd1234_ABCD"
	r.OnOff = "ON"
	err = r.validate()
	assert.EqualError(t, err, "on_off must be set 'on' or 'off'")
	r.OnOff = "on"
	err = r.validate()
	assert.Nil(t, err)
	r.OnOff = "off"
	err = r.validate()
	assert.Nil(t, err)

	// weekday
	// スライス内 0 - 6
	r.Weekday = []int{1}
	err = r.validate()
	assert.Nil(t, err)
	r.Weekday = []int{0, 1, 2, 3, 4, 5, 6}
	err = r.validate()
	assert.Nil(t, err)
	r.Weekday = []int{7}
	err = r.validate()
	assert.EqualError(t, err, "weekday must be between 0 and 6")
	r.Weekday = []int{0, 1, 2, 3, 4, 5, 6, 1}
	err = r.validate()
	assert.EqualError(t, err, "weekday exceeds 7 days")
	r.Weekday = []int{}
	err = r.validate()
	assert.EqualError(t, err, "weekday needs to set at least 1 day")
	r.Weekday = []int{1, 1}
	err = r.validate()
	assert.EqualError(t, err, "weekday duplicates")

	// time
	r.Weekday = []int{1}
	r.Time = ""
	err = r.validate()
	assert.EqualError(t, err, "time must not be empty")
	r.Time = "23:59"
	err = r.validate()
	assert.Nil(t, err)
	r.Time = "0000"
	err = r.validate()
	assert.EqualError(t, err, "invalid time value 0000, must be between 00:00 - 23:59")
	r.Time = "00:60"
	err = r.validate()
	assert.EqualError(t, err, "invalid time value 00:60, must be between 00:00 - 23:59")
	r.Time = "24:00"
	err = r.validate()
	assert.EqualError(t, err, "invalid time value 24:00, must be between 00:00 - 23:59")

	// scriptId
	r.Time = "00:00"
	r.ScriptId = "d1a0e190acc9f55ef7032dedd1a0eeaf98591ac71"
	err = r.validate()
	assert.EqualError(t, err, "script_id must be 40 digits of lower cased hexadecimal")
	r.ScriptId = "D1A0E190ACC9F55EF7032DEDD1A0EEAF98591AC7"
	err = r.validate()
	assert.EqualError(t, err, "script_id must be 40 digits of lower cased hexadecimal")
	r.ScriptId = "g1a0e190acc9f55ef7032dedd1a0eeaf98591ac7"
	err = r.validate()
	assert.EqualError(t, err, "script_id must be 40 digits of lower cased hexadecimal")
}
