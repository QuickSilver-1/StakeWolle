package db

import (
	"context"
	"fmt"

	"referal/internal/config"
	"referal/pkg/log"

	"github.com/go-redis/redis/v8"
)

var(
    RedisConn = NewConnectRedis(*config.AppConfig.PgHost + ":6379", "1234")
)

// NewConnectRedis создает новое подключение к Redis
func NewConnectRedis(addr string, pass string) *redis.Client {
    r := redis.NewClient(&redis.Options{
        Addr: addr,
        Password: pass,
        DB: 0,
    })

    log.Logger.Info("Успешное подключение к Redis")

    return r
}

// NewKey создает новый ключ в Redis
func NewKey(ctx context.Context, key string, value interface{}, out chan string) {
    defer close(out)
    log.Logger.Info(fmt.Sprintf("Создание нового ключа в Redis: %s", key))
    ans := RedisConn.Set(ctx, key, value, 0)
    res, err := ans.Result()

    if err != nil {
        out<- res
        log.Logger.Error(fmt.Sprintf("Ошибка при создании ключа в Redis: %v", err))
        return
    }

    out<- "success"
    log.Logger.Info("Ключ успешно создан в Redis")
}

// GetKey получает значение ключа из Redis
func GetKey(ctx context.Context, key string, out chan string) {
    defer close(out)
    log.Logger.Info(fmt.Sprintf("Получение значения ключа из Redis: %s", key))

    code := RedisConn.Get(ctx, key)
    out<- code.Val()

    log.Logger.Info("Значение ключа успешно получено из Redis")
}

// DelKey удаляет ключ из Redis
func DelKey(ctx context.Context, key string, out chan string) {
    defer close(out)
    log.Logger.Info(fmt.Sprintf("Удаление ключа из Redis: %s", key))
    err := RedisConn.Del(ctx, key).Err()

    if err != nil {
        out<- err.Error()
        log.Logger.Error(fmt.Sprintf("Ошибка при удалении ключа из Redis: %v", err))
        return
    }

    out<- "success"
    log.Logger.Info("Ключ успешно удален из Redis")
}
