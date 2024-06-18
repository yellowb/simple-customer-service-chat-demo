package stream_server

import (
	"log"
	"sync"
)

type ClientChan chan string

var (
	once   sync.Once
	server *SSEventStreamServer
)

// SSEventStreamServer SSE推流服务器
type SSEventStreamServer struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

func GetSSEventStreamServer() *SSEventStreamServer {
	once.Do(func() {
		server = &SSEventStreamServer{
			Message:       make(chan string),
			NewClients:    make(chan chan string),
			ClosedClients: make(chan chan string),
			TotalClients:  make(map[chan string]bool),
		}
		server.run() // 创建时即启动监听
	})
	return server
}

func (s *SSEventStreamServer) run() {
	go s.listen()
}

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
