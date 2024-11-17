package db

import (
	"context"
	"fmt"
	"referal/internal/config"

	"github.com/go-redis/redis/v8"
)

var(
	RedisConn = NewConnectRedis(*config.AppConfig.PgHost + ":6379", "1234")
)

func NewConnectRedis(addr string, pass string) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: pass,
		DB: 0,
	})

	return r
}

func NewKey(ctx context.Context, key string, value interface{}, out chan string) {
	defer close(out)
	err := RedisConn.Set(ctx, key, value, 0)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		out<- err.Err().Error()
		return
	}

	out<- "success"
}

func KeyExist(ctx context.Context, key string, out chan string) {
	defer close(out)
	exist, err := RedisConn.Exists(ctx, key).Result()

	if err != nil {
		out<- err.Error()
		return
	}

	if exist > 0 {
		out<- "exist"
		return
	}

	code := RedisConn.Get(ctx, key)
	out<- code.String()
}

func DelKey(ctx context.Context, key string, out chan string) {
	defer close(out)
	err := RedisConn.Del(ctx, key).Err()

	if err != nil {
		out<- err.Error()
		return
	}

	out<- "success"
}