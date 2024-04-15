package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"weixin_LLM/init/config"
)

var RedisClient *redis.Client
var Ctx context.Context

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.RedisConfigure.Host, config.Config.RedisConfigure.Port), // Redis服务器地址
		Password: "",                                                                                         // 密码，没有则留空
		DB:       0,                                                                                          // 使用默认DB
	})

	ctx := context.Background()
	Ctx = ctx
	// 测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
