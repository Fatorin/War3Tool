package utils

import (
	"errors"
	"strings"
	"unicode"
)

const (
	W3xExtension = ".w3x"
	W3xMime      = "application/octet-stream"
	W3xMapLimit  = 150 << 20 // 150 MB
)

// IsEmptyOrNil 檢查字符串是否為空
func IsEmptyOrNil(str string) bool {
	return len(str) == 0 || str == ""
}

// CheckFileStatus 檢查上傳檔案的狀態
func CheckFileStatus(filename string, mimeType string) error {
	// 檔名不可包含空格
	for _, text := range filename {
		if unicode.IsSpace(text) {
			return errors.New("file name has problem, please remove space and special characters")
		}
	}

	// 檢查副檔名
	if !strings.HasSuffix(filename, W3xExtension) {
		return errors.New("only support .w3x file")
	}

	// 檢查 MIME 類型
	if mimeType != W3xMime {
		return errors.New("not support file")
	}

	return nil
}
