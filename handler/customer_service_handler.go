package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"yellowb.com/chat-demo/db/fake_db"
	"yellowb.com/chat-demo/dto"
	"yellowb.com/chat-demo/stream_server"
)

type CustomerServiceHandler struct {
	// 推流服务器，应当已经运行起来了
	streamServer *stream_server.SSEventStreamServer

	db *fake_db.FakeDb
}

func NewCustomerServiceHandler(streamServer *stream_server.SSEventStreamServer) *CustomerServiceHandler {
	return &CustomerServiceHandler{
		streamServer: streamServer,
		db:           fake_db.GetFakeDb(),
	}
}

// AddMessage 提交一个消息
func (h *CustomerServiceHandler) AddMessage(ctx *gin.Context) {
	msg := new(dto.Message)
	_ = ctx.ShouldBindJSON(msg)
	msg.Ts = time.Now().Unix()

	// 单纯存入db
	h.db.SaveCustomerMessage(msg.User, msg)

	// 广播消息
	h.streamServer.PublishToRedis("hey")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// GetMessagesByUser 获取特定客户的所有对话记录
func (h *CustomerServiceHandler) GetMessagesByUser(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	msgList := h.db.GetCustomerMessages(userId)
	if msgList == nil {
		// 找不到，返回空对象
		msgList = new(dto.CustomerMessages)
	}
	ctx.JSON(http.StatusOK, msgList)
}
