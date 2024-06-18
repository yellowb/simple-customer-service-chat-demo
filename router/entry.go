package router

import "github.com/gin-gonic/gin"

func InitRouter() (*gin.Engine, error) {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 静态html文件
	r.StaticFile("/", "./public/index.html")

	rg := r.Group("/api")
	RegisterCustomerServiceRouter(rg)

	return r, nil
}
