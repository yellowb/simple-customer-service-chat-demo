package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"sync/atomic"
	"yellowb.com/chat-demo/stream_server"
)

// SSEventStreamHandler SSEvent推流接口
type SSEventStreamHandler struct {
	// 推流服务器，应当已经运行起来了
	streamServer *stream_server.SSEventStreamServer
	// 用来生成客户端ID
	counter atomic.Uint64
}

func NewSSEventStreamHandler(streamServer *stream_server.SSEventStreamServer) *SSEventStreamHandler {
	return &SSEventStreamHandler{
		streamServer: streamServer,
	}
}

// AcceptNewClient 接受一个新的客户端连上来
func (h *SSEventStreamHandler) AcceptNewClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 给每一个客户端链接一个ID，当前没啥用，可以放在ctx中后续取出
		currentCount := h.counter.Add(1)
		streamId := fmt.Sprintf("StreamID-%d", currentCount)

		// 代表一个客户端链接，后面往客户端推数据要靠它
		clientChan := make(stream_server.ClientChan)

		// 加到SSE Server的客户端集合中
		h.streamServer.NewClients <- clientChan
		fmt.Printf("client - %s is added\n", streamId)

		// 客户端主动断开链接时会走到这里
		defer func() {
			fmt.Printf("client - %s is lost\n", streamId)
			// 把客户端链接从SSE Server中移除
			h.streamServer.ClosedClients <- clientChan
		}()

		// 只要客户端不断开链接，下面这个函数会一直循环执行
		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent(msg.Type, msg.Body)
				return true
			}
			return false
		})
	}
}
