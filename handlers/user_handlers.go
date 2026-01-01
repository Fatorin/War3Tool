package handlers

import (
	"net/http"
	"os"
	"war3tool/services"

	"github.com/gin-gonic/gin"
)

var usersFolderPath = ""
var templatesPath = ""

// InitUserHandlers 初始化用戶相關的處理器
func InitUserHandlers(usersPath, tmplPath string) {
	usersFolderPath = usersPath
	templatesPath = tmplPath
}

// RegisterPage 渲染註冊頁面
func RegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", nil)
}

// RegisterUser 註冊新用戶
func RegisterUser(ctx *gin.Context) {
	var userData services.User

	if err := ctx.ShouldBindJSON(&userData); err != nil {
		println(err.Error())
		ctx.JSON(
			http.StatusBadRequest, gin.H{"msg": "系統出現異常，請檢查異常。"})
		return
	}

	// 驗證驗證碼
	valid := os.Getenv("FA_VALID")
	if userData.Valid != valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "驗證碼不正確，請確認驗證碼。"})
		return
	}

	// 創建用戶
	pwd, err := services.CreateUser(userData, usersFolderPath, templatesPath, "")
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "系統出現異常，請檢查異常。"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"username": userData.Username,
			"password": pwd,
		})
}
