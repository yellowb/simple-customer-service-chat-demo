package router

import (
	"github.com/gin-gonic/gin"
	"yellowb.com/chat-demo/handler"
	"yellowb.com/chat-demo/middleware"
	"yellowb.com/chat-demo/stream_server"
)

func RegisterCustomerServiceRouter(rg *gin.RouterGroup) {
	sseHandler := handler.NewSSEventStreamHandler(stream_server.GetSSEventStreamServer())
	csHandler := handler.NewCustomerServiceHandler(stream_server.GetSSEventStreamServer())

	rg.GET("/stream", middleware.AddSSEventHeader(), sseHandler.AcceptNewClient())
	rg.POST("/add_message", csHandler.AddMessage)
	rg.GET("/messages", csHandler.GetMessagesByUser)

}
