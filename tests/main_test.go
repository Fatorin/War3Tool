package main

import (
	"testing"
	"war3tool/utils"
)

func TestCheckFileStatus(t *testing.T) {
	err := utils.CheckFileStatus("fate.w3x", "application/octet-stream")

	if err != nil {
		t.Error(err)
	}
}

func TestLoadMapsFolder(t *testing.T) {
	// 注意：需要有效的地圖資料夾路徑
	// 這個測試可能需要設置環境或使用 mock
	// mapsDict, err := handlers.LoadMapsFolder("./data/maps")
	// if err != nil {
	// 	t.Error(err)
	// }
}
