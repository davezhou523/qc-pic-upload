package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	RedisClient struct {
		Addr     string
		Password string
		DB       int
	}
}
