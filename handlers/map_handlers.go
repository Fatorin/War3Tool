package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"war3tool/services"
	"war3tool/utils"

	"github.com/gin-gonic/gin"
)

var (
	mapsNameDict   = make(map[string]string)
	mapsFolderPath = ""
	AnalysisW3x    func(string) (string, error)
)

// InitMapHandlers 初始化地圖相關的處理器
func InitMapHandlers(mapsDictRef map[string]string, mapsPath string, analysisFunc func(string) (string, error)) {
	mapsNameDict = mapsDictRef
	mapsFolderPath = mapsPath
	AnalysisW3x = analysisFunc
}

// LoadMapsFolder 從資料夾加載所有地圖
func LoadMapsFolder(mapsPath string) (map[string]string, error) {
	mapsDict := make(map[string]string)

	files, err := os.ReadDir(mapsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		mapPath := filepath.Join(mapsPath, file.Name())
		mapAnalysisName, err := services.AnalysisW3x(mapPath)

		if err != nil {
			panic(err)
		}

		mapsDict[file.Name()] = mapAnalysisName
	}

	return mapsDict, nil
}

// UploadPage 渲染上傳頁面
func UploadPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "upload.html", nil)
}

// GetFiles 獲取所有地圖檔案列表
func GetFiles(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"files": mapsNameDict,
	})
}

// UploadSingleFile 上傳單個地圖檔案
func UploadSingleFile(ctx *gin.Context) {
	const w3xMapLimitSize = utils.W3xMapLimit

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(w3xMapLimitSize))
	if len(ctx.Errors) > 0 {
		return
	}

	// 驗證驗證碼
	valid := ctx.PostForm("valid")
	validCode := os.Getenv("FA_VALID")
	if valid != validCode {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "驗證碼不正確，請確認驗證碼。",
		})
		return
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	// 檢查檔案狀態
	err = utils.CheckFileStatus(file.Filename, file.Header.Get("Content-type"))
	if err != nil {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	// 檢查檔案名稱是否重複
	if _, ok := mapsNameDict[file.Filename]; ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：檔案名稱重複"),
		})
		return
	}

	// 保存上傳的檔案
	savePath := filepath.Join(mapsFolderPath, file.Filename)

	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	// 分析地圖檔案
	realName, err := AnalysisW3x(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintln("解析錯誤，原因：", err.Error()),
		})
		return
	}

	// 將地圖添加到字典
	mapsNameDict[file.Filename] = realName

	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintln("檔名：", file.Filename, "圖名：", realName),
	})
}
