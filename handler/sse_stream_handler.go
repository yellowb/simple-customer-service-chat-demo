package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"sync/atomic"
	"yellowb.com/chat-demo/stream_server"
)

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
		currentCount := h.counter.Load()
		streamId := fmt.Sprintf("StreamID-%d", currentCount)
		h.counter.Add(1)

		// Initialize client channel
		clientChan := make(stream_server.ClientChan)

		// Send new connection to event server
		h.streamServer.NewClients <- clientChan
		fmt.Printf("client - %s is added\n", streamId)

		defer func() {
			fmt.Printf("client - %s is lost\n", streamId)
			// Send closed connection to event server
			h.streamServer.ClosedClients <- clientChan
		}()

		c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientChan; ok {
				c.SSEvent("message", fmt.Sprintf("%s - %s", msg, streamId))
				return true
			}
			return false
		})
	}
}
