package main

import (
	"bytes"
	"errors"
	"fmt"
	"go_system/pvpgnhash"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const w3xExtension = ".w3x"
const w3xMime = "application/octet-stream"
const w3xMapLimitSize = 150 << 20

var mapsNameDict = make(map[string]string)
var usersFolderPath = ""
var mapsFolderPath = ""
var currentDir = ""

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		println("找不到執行環境資料夾")
		return
	}

	valid := os.Getenv("FA_VALID")
	if isEmptyOrNil(valid) {
		println("驗證碼為空，請加入環境變數")
		return
	}

	usersFolderPath = os.Getenv("USERS_FOLDER_PATH")
	if isEmptyOrNil(usersFolderPath) {
		println("用戶資料夾參數為空，請加入環境變數")
		return
	}

	usersFolderPath = filepath.Join(currentDir, usersFolderPath)
	if err := checkFolder(usersFolderPath); err != nil {
		println("找不到用戶資料夾，請加入環境變數。路徑：", usersFolderPath)
		return
	}

	mapsFolderPath = os.Getenv("MAPS_FOLDER_PATH")
	if isEmptyOrNil(mapsFolderPath) {
		println("地圖資料夾參數為空，請加入環境變數")
		return
	}

	mapsFolderPath = filepath.Join(currentDir, mapsFolderPath)
	if err := checkFolder(mapsFolderPath); err != nil {

		println("找不到地圖資料夾，請確認路徑是否正確。路徑：", mapsFolderPath)
		return
	}

	LoadMapsFolder()
	GinInit()
}

func LoadMapsFolder() {

	for k := range mapsNameDict {
		delete(mapsNameDict, k)
	}

	files, err := ioutil.ReadDir(mapsFolderPath)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		mapPath := filepath.Join(mapsFolderPath, file.Name())
		mapAnalysisName, err := AnalysisW3x(mapPath)

		if err != nil {
			panic(err)
		}

		mapsNameDict[file.Name()] = mapAnalysisName
	}

}

func AnalysisW3x(path string) (mapName string, err error) {
	file, err := os.Open(path)
	mapName = ""

	err = nil

	if err != nil {
		err = errors.New("not found file")
	}

	defer file.Close()

	buffer := make([]byte, 128)

	_, err = file.Read(buffer)
	if err != nil {
		err = errors.New("read has problem")
	}

	//處理後面不要的文字
	newBuffer := buffer[8:]
	emptyByte := []byte{0x00}
	pos := bytes.Index(newBuffer, emptyByte)
	mapName = string(newBuffer[:pos])
	mapName = strings.ReplaceAll(mapName, "|r", "")
	//處理多餘的|c
	for {
		pos = strings.Index(mapName, "|c")
		if pos == -1 {
			break
		}
		mapName = strings.ReplaceAll(mapName, mapName[pos:pos+10], "")
	}

	return
}

func GinInit() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200", "http://localhost:80", "http://localhost:8080", "http://220.133.54.223:8080"}
	r.Use((cors.New(config)))
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.POST("/UploadFile", UploadSingleFile)
	r.GET("/GetFiles", GetFiles)
	r.POST("/Register", RegisterUser)
	r.LoadHTMLGlob("templates/*")
	r.GET("/signup", RegisterPage)
	r.GET("/upload", UploadPage)
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}

func UploadPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "upload.html", nil)
}

func GetFiles(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"files": mapsNameDict,
	})
}

type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	Valid    string `json:"valid"`
}

func RegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

func RegisterUser(ctx *gin.Context) {
	var data User

	if err := ctx.ShouldBindJSON(&data); err != nil {
		println(err.Error())
		ctx.JSON(
			http.StatusBadRequest, gin.H{"msg": "系統出現異常，請檢查異常。"})
		return
	}

	valid := os.Getenv("FA_VALID")
	if data.Valid != valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "驗證碼不正確，請確認驗證碼。"})
		return
	}

	templateFile := "user_template"
	if data.IsAdmin {
		templateFile = "admin_template"
	}

	userFilePath := filepath.Join(usersFolderPath, data.Username)
	if fileExists(userFilePath) {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "此用戶名已存在"})
		return
	}

	count, err := countFilesInFolder(usersFolderPath)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "系統出現異常，請檢查異常。"})
		return
	}

	pwd := "fate" + generateRandomPassword(6)
	hash := pvpgnhash.GetHash(pwd)
	filefolder := filepath.Join(currentDir, "files")
	filedata, err := readFile(filefolder, templateFile)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "系統出現異常，請檢查異常。"})
		return
	}

	filedata = strings.Replace(filedata, "{{userid}}", strconv.Itoa(count), -1)
	filedata = strings.Replace(filedata, "{{username}}", data.Username, -1)
	filedata = strings.Replace(filedata, "{{passhash1}}", hash, -1)

	filepath := filepath.Join(usersFolderPath, data.Username)
	err = saveToFile(filedata, filepath)
	if err != nil {
		println(err.Error())
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg": "系統出現異常，請檢查異常。",
			})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"username": data.Username,
			"password": pwd,
		})
}

func UploadSingleFile(ctx *gin.Context) {

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, w3xMapLimitSize)
	if len(ctx.Errors) > 0 {
		return
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	err = CheckFileStatus(file.Filename, file.Header.Get("Content-type"))
	if err != nil {
		ctx.JSON(http.StatusUnsupportedMediaType, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	if _, ok := mapsNameDict[file.Filename]; ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：檔案名稱重複"),
		})
		return
	}

	savePath := filepath.Join(mapsFolderPath, file.Filename)

	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintln("上傳失敗，原因：", err.Error()),
		})
		return
	}

	realName, err := AnalysisW3x(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintln("解析錯誤，原因：", err.Error()),
		})
		return
	}

	mapsNameDict[file.Filename] = realName

	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintln("檔名：", file.Filename, "圖名：", realName),
	})
}

func CheckFileStatus(filename string, mimeType string) error {
	//檔名不可空白
	for _, text := range filename {
		if unicode.IsSpace(text) {
			return errors.New("file name has problem, please remove space and special characters")
		}
	}

	if !strings.HasSuffix(filename, w3xExtension) {
		return errors.New("only support .w3x file")
	}

	//檢查是否為application/octet-stream
	if mimeType != w3xMime {
		return errors.New("not support file")
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func checkFolder(folderPath string) error {
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("資料夾不存在：%s", folderPath)
		}
		return err
	}

	return nil
}

func generateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())

	password := ""
	for i := 0; i < length; i++ {
		password += fmt.Sprintf("%d", rand.Intn(10))
	}

	return password
}

func countFilesInFolder(folderPath string) (int, error) {
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

func readFile(folderPath, fileName string) (string, error) {
	filePath := folderPath + "/" + fileName
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("讀取失敗：%s", filePath)
	}

	return string(fileData), nil
}

func saveToFile(text, filePath string) error {
	err := ioutil.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

func isEmptyOrNil(str string) bool {
	return len(str) == 0 || str == ""
}
