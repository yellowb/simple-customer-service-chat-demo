package stream_server

import (
	"log"
	"sync"
	"time"
)

const (
	heartbeatInterval = time.Second * 10 // 心跳保活间隔
)

// SSEvent SSE事件消息体
type SSEvent struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

// ClientChan 面向单个客户端的推流Channel
type ClientChan chan *SSEvent

var (
	once   sync.Once
	server *SSEventStreamServer
)

// SSEventStreamServer SSE推流服务器
type SSEventStreamServer struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message ClientChan

	// New client connections
	NewClients chan ClientChan

	// Closed client connections
	ClosedClients chan ClientChan

	// Total client connections
	TotalClients map[ClientChan]bool
}

func GetSSEventStreamServer() *SSEventStreamServer {
	once.Do(func() {
		server = &SSEventStreamServer{
			Message:       make(chan *SSEvent),
			NewClients:    make(chan ClientChan),
			ClosedClients: make(chan ClientChan),
			TotalClients:  make(map[ClientChan]bool),
		}
		server.run() // 创建时即启动监听
	})
	return server
}

func (s *SSEventStreamServer) run() {
	go s.listen()
	go s.keepAlive()
}

// 监听新客户端链接，并把每一条消息广播到所有客户端
func (s *SSEventStreamServer) listen() {
	for {
		select {
		// Add new available client
		case client := <-s.NewClients:
			s.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(s.TotalClients))

		// Remove closed client
		case client := <-s.ClosedClients:
			delete(s.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(s.TotalClients))

		// Broadcast message to client
		case eventMsg := <-s.Message:
			for clientMessageChan := range s.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

// 定期广播心跳消息到所有客户端以保证链接存活
func (s *SSEventStreamServer) keepAlive() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		s.Message <- &SSEvent{
			Type: "heartbeat",
		}
	}
}
