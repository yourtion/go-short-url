package redis

import (
	"github.com/go-redis/redis"

	"github.com/yourtion/go-short-url/internal/base/config"
	"github.com/yourtion/go-short-url/internal/base/logger"
)

var log *logger.Entry

// Redis Client
var Client *redis.Client

// 打开 Redis 连接
func Open(opts *config.RedisConfig) *redis.Client {
	log = logger.NewModuleLogger("redis")
	conf := redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
		PoolSize: opts.PoolSize,
		OnConnect: func(conn *redis.Conn) error {
			log.Debugln("Client Connected")
			return nil
		},
	}
	Client = redis.NewClient(&conf)
	return Client
}

// 获取 Redis Key
func GetKey(key string) string {
	if config.Config.Redis.Prefix == "" {
		return key
	}
	return config.Config.Redis.Prefix + ":" + key
}
