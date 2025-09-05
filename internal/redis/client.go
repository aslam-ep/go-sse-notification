package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	RedisAddressPrefix = "redis://"

	NotificationTopic = "notificaions"
)

type Client struct {
	RDB *redis.Client
}

func NewClient(address, password string, db int) *Client {
	var rdb *redis.Client

	if strings.HasPrefix(address, RedisAddressPrefix) {
		opt, _ := redis.ParseURL(address)
		rdb = redis.NewClient(opt)
	} else {
		rdb = redis.NewClient(&redis.Options{Addr: address, Password: password, DB: db})
	}

	return &Client{RDB: rdb}
}

func (c *Client) Publish(ctx context.Context, channel, message string) error {
	return c.RDB.Publish(ctx, channel, message).Err()
}

func (c *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.RDB.Subscribe(ctx, channel)
}
