package stream_server

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
	"yellowb.com/chat-demo/client/redis_client"
	"yellowb.com/chat-demo/constants"
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
	// 心跳消息入口
	Heartbeat chan *SSEvent

	// Events are pushed to this channel by the main events-gathering routine
	Message chan *SSEvent

	// New client connections
	NewClients chan ClientChan

	// Closed client connections
	ClosedClients chan ClientChan

	// Total client connections
	TotalClients map[ClientChan]bool

	// Redis client
	redisClient *redis_client.Client
}

func GetSSEventStreamServer() *SSEventStreamServer {
	once.Do(func() {
		server = &SSEventStreamServer{
			Heartbeat:     make(chan *SSEvent),
			Message:       make(chan *SSEvent),
			NewClients:    make(chan ClientChan),
			ClosedClients: make(chan ClientChan),
			TotalClients:  make(map[ClientChan]bool),
			redisClient:   redis_client.GetClient(),
		}
		server.run() // 创建时即启动监听
	})
	return server
}

func (s *SSEventStreamServer) run() {
	go s.listen()
	go s.keepAlive()
	go s.subscribeToRedis()
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

			// 广播心跳信号
		case heartbeatMsg := <-s.Heartbeat:
			for clientMessageChan := range s.TotalClients {
				clientMessageChan <- heartbeatMsg
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
		s.Heartbeat <- &SSEvent{
			Type: "heartbeat",
		}
	}
}

func (s *SSEventStreamServer) subscribeToRedis() {
	subscriber := s.redisClient.Client.Subscribe(context.Background(), constants.CustomerServiceMsgNotifyChan)
	defer subscriber.Close()

	ch := subscriber.Channel(redis.WithChannelSize(1000), redis.WithChannelSendTimeout(5*time.Second))

	// 如果Redis Chan中没有数据，这里会被堵塞
	for msg := range ch {
		s.Message <- &SSEvent{
			Type: "message",
			Body: msg.Payload,
		}
	}
}

func (s *SSEventStreamServer) PublishToRedis(payload string) error {
	err := s.redisClient.Client.Publish(context.Background(), constants.CustomerServiceMsgNotifyChan, payload).Err()
	return err
}
