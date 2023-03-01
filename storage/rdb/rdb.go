package rdb

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

// Connection configuration
type Option struct {
	Addr     string // example => 127.0.0.1:6379
	Password string
	DB       int
}

// Storage structure
type Storage struct {
	Client *redis.Client
}

var ctx = context.Background()

// Create Redis-V8 storage instance
func NewStorage(option *Option) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
		DB:       option.DB,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("Ping failed")
	}
	return &Storage{Client: client}, nil
}
