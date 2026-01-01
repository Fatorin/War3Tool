package main

import (
	"fmt"
	"os"
	"path/filepath"
	"war3tool/handlers"
	"war3tool/utils"
)

var currentDir = ""

func main() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		println("找不到執行環境資料夾")
		return
	}

	// 驗證碼檢查
	valid := os.Getenv("FA_VALID")
	if utils.IsEmptyOrNil(valid) {
		println("驗證碼為空，請加入環境變數")
		return
	}

	// 用戶資料夾檢查
	usersFolderPath := os.Getenv("USERS_FOLDER_PATH")
	if utils.IsEmptyOrNil(usersFolderPath) {
		println("用戶資料夾參數為空，請加入環境變數")
		return
	}

	usersFolderPath = filepath.Join(currentDir, usersFolderPath)
	if err := utils.CheckFolder(usersFolderPath); err != nil {
		println("找不到用戶資料夾，請加入環境變數。路徑：", usersFolderPath)
		return
	}

	// 地圖資料夾檢查
	mapsFolderPath := os.Getenv("MAPS_FOLDER_PATH")
	if utils.IsEmptyOrNil(mapsFolderPath) {
		println("地圖資料夾參數為空，請加入環境變數")
		return
	}

	mapsFolderPath = filepath.Join(currentDir, mapsFolderPath)
	if err := utils.CheckFolder(mapsFolderPath); err != nil {
		println("找不到地圖資料夾，請確認路徑是否正確。路徑：", mapsFolderPath)
		return
	}

	// 加載地圖資料夾
	mapsDict, err := handlers.LoadMapsFolder(mapsFolderPath)
	if err != nil {
		println("加載地圖資料夾失敗：", err.Error())
		return
	}

	// 模板資料夾路徑
	templatesPath := filepath.Join(currentDir, "data/templates")
	if err := utils.CheckFolder(templatesPath); err != nil {
		println("找不到模板資料夾，請確認路徑是否正確。路徑：", templatesPath)
		return
	}

	// 初始化路由並運行
	r := handlers.InitRouter(mapsDict, mapsFolderPath, usersFolderPath, templatesPath)
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
