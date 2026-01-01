package services

import (
	"path/filepath"
	"strconv"
	"strings"
	"war3tool/crypto"
	"war3tool/utils"
)

// User 用戶數據結構
type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	Valid    string `json:"valid"`
}

// CreateUser 創建新用戶
func CreateUser(userData User, usersFolderPath, templatesPath, currentDir string) (string, error) {
	// 檢查用戶是否已存在
	userFilePath := filepath.Join(usersFolderPath, userData.Username)
	if utils.FileExists(userFilePath) {
		return "", utils.UserAlreadyExistsError()
	}

	// 計算用戶 ID
	count, err := utils.CountFilesInFolder(usersFolderPath)
	if err != nil {
		return "", err
	}

	// 生成隨機密碼並加密
	pwd := "fate" + utils.GenerateRandomPassword(6)
	hash := crypto.GetHash(pwd)

	// 讀取模板檔案
	templateFile := "user_template"
	if userData.IsAdmin {
		templateFile = "admin_template"
	}

	filedata, err := utils.ReadFile(templatesPath, templateFile)
	if err != nil {
		return "", err
	}

	filedata = strings.ReplaceAll(filedata, "{{userid}}", strconv.Itoa(count))
	filedata = strings.ReplaceAll(filedata, "{{username}}", userData.Username)
	filedata = strings.ReplaceAll(filedata, "{{passhash1}}", hash)

	// 保存用戶檔案
	userFilePath = filepath.Join(usersFolderPath, userData.Username)
	err = utils.SaveToFile(filedata, userFilePath)
	if err != nil {
		return "", err
	}

	return pwd, nil
}
