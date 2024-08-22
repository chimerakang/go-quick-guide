package routes

import (
	"go-todo/controllers"
	"go-todo/controllers/jwt"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	// 創建 Gin 引擎
	r := gin.Default()

	// 加載 HTML 模板
	r.LoadHTMLGlob("views/*")

	// 設置靜態文件路徑（如果有的話）
	r.Static("/static", "./static")

	// 公開路由
	r.GET("/login", controllers.ShowLogin)
	r.POST("/login", controllers.ProcessLogin)
	r.GET("/logout", controllers.ProcessLogout)

	// 受保護的路由
	protected := r.Group("/")
	protected.Use(jwt.AuthMiddleware())
	{
		protected.GET("/", controllers.Show)
		protected.POST("/add", controllers.Add)
		protected.GET("/delete/:id", controllers.Delete)
		protected.GET("/complete/:id", controllers.Complete)
	}

	// // 設置路由
	// r.GET("/", controllers.Show)
	// r.POST("/add", controllers.Add)
	// r.GET("/delete/:id", controllers.Delete)
	// r.GET("/complete/:id", controllers.Complete)
	// r.GET("/login", controllers.ShowLogin)
	// r.POST("/login", controllers.ProcessLogin)
	// r.GET("/logout", controllers.ProcessLogout)

	return r
}
