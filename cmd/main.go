package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yellowb.com/chat-demo/router"
	"yellowb.com/chat-demo/stream_server"
)

func main() {
	// start api server...
	r, err := router.InitRouter()
	if err != nil {
		panic(fmt.Sprintf("[main] init router error: %s", err.Error()))
	}
	webServer := &http.Server{
		Addr:           ":8085",
		Handler:        r,
		ReadTimeout:    30 * time.Minute,
		WriteTimeout:   30 * time.Minute,
		IdleTimeout:    30 * time.Minute,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	go func() {
		// 启动WebServer开始监听请求
		if err := webServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[main] start web server error: %v", err)
		}
	}()

	// TODO: 后面删掉
	go func() {
		streamServer := stream_server.GetSSEventStreamServer()
		for {
			time.Sleep(time.Second * 10)
			now := time.Now().Format("2006-01-02 15:04:05")
			currentTime := fmt.Sprintf("The Current Time Is %v", now)

			// Send current time to clients message channel
			streamServer.Message <- currentTime
		}
	}()

	// graceful shutdown...
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] shutdown web server ...")

	// 关闭WebServer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := webServer.Shutdown(ctx); err != nil {
		log.Fatal("[main] server shutdown: ", err)
	}
	log.Println("[main] shutdown web server ok")

	log.Println("[main] server exiting")
}
