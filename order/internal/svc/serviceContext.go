package svc

import (
	"github.com/redis/go-redis/v9"
	"qc/order/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     c.RedisClient.Addr,
			Password: c.RedisClient.Password,
			DB:       c.RedisClient.DB,
		}),
	}
}
