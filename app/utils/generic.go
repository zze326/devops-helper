package utils

import "unicode/utf8"

// FilterUTF8 过滤 byte 中非 utf8 字符
func FilterUTF8(b []byte) []byte {
	var result []byte
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			result = append(result, []byte("\ufffd")...)
		} else {
			result = append(result, b[:size]...)
		}
		b = b[size:]
	}
	return result
}

func InSlice[T comparable](dest T, array []T) bool {
	for _, v := range array {
		if v == dest {
			return true
		}
	}
	return false
}
