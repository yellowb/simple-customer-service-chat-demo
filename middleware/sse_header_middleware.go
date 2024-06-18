package middleware

import "github.com/gin-gonic/gin"

// AddSSEventHeader 给浏览器的响应中添加SSE header，让浏览器保持链接不断
func AddSSEventHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		// 禁止网关缓存SSE流（https://blog.csdn.net/u013534071/article/details/131500873）
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Next()
	}
}
