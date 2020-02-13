package utils

// 引数の bool が true のとき1を、そうでないとき0を返す。
func BoolToInt(b bool) int {
	ret := 0
	if b {
		ret = 1
	}

	return ret
}

// 引数の bool が true のとき「on」、そうでないとき「off」を返す。
func BoolToOnOff(b bool) string {
	ret := "off"
	if b {
		ret = "on"
	}

	return ret
}

// 引数が on のとき true を、そうでないとき false を返す。
func OnOffToBool(onOff string) bool {
	return onOff == "on"
}

// 引数が 1 のとき true を、そうでないとき false を返す。
func IntToBool(i int) bool {
	return i == 1
}
