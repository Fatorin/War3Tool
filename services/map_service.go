package services

import (
	"bytes"
	"errors"
	"strings"
	"war3tool/utils"
)

// AnalysisW3x 分析 W3x 檔案並提取地圖名稱
func AnalysisW3x(path string) (mapName string, err error) {
	file, err := utils.OpenFile(path)
	mapName = ""

	if err != nil {
		err = errors.New("not found file")
		return
	}

	defer file.Close()

	buffer := make([]byte, 128)

	_, err = file.Read(buffer)
	if err != nil {
		err = errors.New("read has problem")
		return
	}

	// 處理後面不要的文字
	newBuffer := buffer[8:]
	emptyByte := []byte{0x00}
	pos := bytes.Index(newBuffer, emptyByte)
	mapName = string(newBuffer[:pos])
	mapName = strings.ReplaceAll(mapName, "|r", "")

	// 處理多餘的 |c
	for {
		pos = strings.Index(mapName, "|c")
		if pos == -1 {
			break
		}
		mapName = strings.ReplaceAll(mapName, mapName[pos:pos+10], "")
	}

	return
}
