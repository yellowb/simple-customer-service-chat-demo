package redis_client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

const (
	addr = "127.0.0.1:6379"
	db   = 0
)

var (
	once   sync.Once
	client *Client
)

type Client struct {
	Client *redis.Client
}

func GetClient() *Client {
	once.Do(func() {
		client, _ = NewClient()
	})
	return client
}

func NewClient() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	// ping一下看看通不通
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: client,
	}, nil
}
