package utils

import (
	"errors"
	"os"
)

// OpenFile 打開檔案
func OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

// UserAlreadyExistsError 用戶已存在錯誤
func UserAlreadyExistsError() error {
	return errors.New("此用戶名已存在")
}
