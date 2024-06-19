package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"yellowb.com/chat-demo/db/fake_db"
	"yellowb.com/chat-demo/dto"
	"yellowb.com/chat-demo/stream_server"
)

// CustomerServiceHandler 客户提问、客服回答接口
type CustomerServiceHandler struct {
	// 推流服务器，应当已经运行起来了
	// PS：正确来说这个handler不应该包含这个东西，当前仅限于Demo
	streamServer *stream_server.SSEventStreamServer
	// DAO
	db *fake_db.FakeDb
}

func NewCustomerServiceHandler(streamServer *stream_server.SSEventStreamServer) *CustomerServiceHandler {
	return &CustomerServiceHandler{
		streamServer: streamServer,
		db:           fake_db.GetFakeDb(),
	}
}

// AddMessage 提交一个消息（客户和客服通用）
func (h *CustomerServiceHandler) AddMessage(ctx *gin.Context) {
	msg := new(dto.Message)
	_ = ctx.ShouldBindJSON(msg)
	msg.Ts = time.Now().Unix()

	// 单纯存入db
	err := h.db.SaveCustomerMessage(msg.User, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 广播消息
	h.streamServer.PublishToRedis("hey")

	// TODO: 如果是客服回复客户，正常来说还要向WOA推送消息

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// GetMessagesByUser 获取特定客户的所有对话记录
func (h *CustomerServiceHandler) GetMessagesByUser(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	msgList, err := h.db.GetCustomerMessages(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if msgList == nil {
		// 找不到，返回空对象
		msgList = new(dto.CustomerMessages)
	}
	ctx.JSON(http.StatusOK, msgList)
}
