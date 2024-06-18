package router

import "github.com/gin-gonic/gin"

func InitRouter() (*gin.Engine, error) {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 静态html文件
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/ask", "./public/ask.html")
	r.StaticFile("/answer", "./public/answer.html")

	rg := r.Group("/api")
	RegisterCustomerServiceRouter(rg)

	return r, nil
}
