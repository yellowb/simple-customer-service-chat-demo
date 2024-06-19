package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yellowb.com/chat-demo/client/redis_client"
	"yellowb.com/chat-demo/router"
)

var (
	port = flag.Int("p", 8085, "The port to listen on")
)

func main() {
	flag.Parse()

	// start api server...
	r, err := router.InitRouter()
	if err != nil {
		panic(fmt.Sprintf("[main] init router error: %s", err.Error()))
	}
	webServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        r,
		ReadTimeout:    30 * time.Minute,
		WriteTimeout:   30 * time.Minute, // 这个可以适当调大，不然每过一段时间客户端就好被断开一次，不过也无大碍
		IdleTimeout:    30 * time.Minute,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	go func() {
		// 启动WebServer开始监听请求
		if err := webServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[main] start web server error: %v", err)
		}
	}()

	// graceful shutdown...
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] shutdown web server ...")

	// 清除redis中的内容
	redis_client.GetClient().Client.FlushDB(context.Background()).Result()
	log.Println("[main] clear redis data ...")

	// 关闭WebServer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := webServer.Shutdown(ctx); err != nil {
		log.Fatal("[main] server shutdown: ", err)
	}
	log.Println("[main] shutdown web server ok")

	log.Println("[main] server exiting")
}
