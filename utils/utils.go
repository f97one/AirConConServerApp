package utils

// 引数の bool が true のとき1を、そうでないとき0を返す。
func BoolToInt(b bool) int {
	ret := 0
	if b {
		ret = 1
	}

	return ret
}
