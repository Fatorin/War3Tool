package handlers

import (
	"war3tool/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由
func InitRouter(mapsDict map[string]string, mapsPath, usersPath, tmplPath string) *gin.Engine {
	r := gin.Default()

	// CORS 配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost",
		"http://localhost:8080",
	}
	r.Use(cors.New(config))

	// 靜態檔案
	r.Use(static.Serve("/", static.LocalFile("./static", true)))

	// 初始化各個處理器
	InitMapHandlers(mapsDict, mapsPath, services.AnalysisW3x)
	InitUserHandlers(usersPath, tmplPath)

	// 路由定義
	r.POST("/UploadFile", UploadSingleFile)
	r.GET("/GetFiles", GetFiles)
	r.POST("/Register", RegisterUser)
	r.LoadHTMLGlob("assets/templates/*")
	r.GET("/signup", RegisterPage)
	r.GET("/upload", UploadPage)

	return r
}
