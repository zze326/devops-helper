package utils

import (
	"os"
	"strings"
)

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

func OpenOrCreateFile(filepath string) (*os.File, error) {
	dirPath := filepath[:strings.LastIndex(filepath, "/")]
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 确保目标文件不存在
func EnsureFileNotExists(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}
	return nil
}
