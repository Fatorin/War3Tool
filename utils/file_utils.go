package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists 檢查檔案是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CheckFolder 檢查資料夾是否存在
func CheckFolder(folderPath string) error {
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("資料夾不存在：%s", folderPath)
		}
		return err
	}
	return nil
}

// CountFilesInFolder 計算資料夾內的檔案數
func CountFilesInFolder(folderPath string) (int, error) {
	fileCount := 0

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileCount++
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}

// ReadFile 讀取檔案內容
func ReadFile(folderPath, fileName string) (string, error) {
	filePath := folderPath + "/" + fileName
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("讀取失敗：%s", filePath)
	}

	return string(fileData), nil
}

// SaveToFile 將內容保存到檔案
func SaveToFile(text, filePath string) error {
	err := os.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}
