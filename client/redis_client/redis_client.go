package redis_client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	addr = "127.0.0.1:6379"
	db   = 0
)

type Client struct {
	client *redis.Client
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
		client: client,
	}, nil
}
