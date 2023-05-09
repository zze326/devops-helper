package utils

import "os"

func FileMode(mode os.FileMode) string {
	var buf [32]byte
	b := buf[:0]
	switch mode & os.ModeType {
	case os.ModeDir:
		b = append(b, 'd')
	case os.ModeSymlink:
		b = append(b, 'l')
	default:
		b = append(b, '-')
	}
	for i, c := range "rwxrwxrwx" {
		if mode&(1<<(8-i)) == 0 {
			c = '-'
		}
		b = append(b, byte(c))
	}
	return string(b)
}

// 创建或打开文件，返回文件句柄
func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
}
