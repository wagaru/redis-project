package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/wagaru/redis-project/domain"
)

type redisRepo struct {
	client *redis.Client
}

type Repo interface {
	domain.PostRepository
	domain.UserRepository
}

func NewRedisRepo(addr string) (Repo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &redisRepo{
		client: client,
	}, nil
}
