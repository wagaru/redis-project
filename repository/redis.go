package repository

import (
	"github.com/go-redis/redis/v8"
)

type RedisRepo struct {
	client *redis.Client
}
