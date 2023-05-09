package utils

import (
	"fmt"
	"strings"
)

// IsEmpty
// desc: 判断字符串是否为空
func IsEmpty(str string) bool {
	return len(strings.TrimSpace(fmt.Sprintf("%v", str))) == 0
}
