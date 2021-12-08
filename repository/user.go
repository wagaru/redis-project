package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/wagaru/redis-project/domain"
)

type UserRepo interface {
	domain.UserRepository
}

func NewRedisUserRepo(addr string) (UserRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisRepo{
		client: client,
	}, nil
}

func (r *RedisRepo) StoreUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		res, err := r.client.Incr(ctx, "user:").Result()
		if err != nil {
			return err
		}
		user.ID = strconv.FormatInt(res, 10)
	}
	err := r.client.SAdd(ctx, "users:", "user:"+user.ID).Err()
	if err != nil {
		return err
	}
	err = r.client.HSet(ctx, "tokens:", user.Token, "user:"+user.ID).Err()
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, "user:"+user.ID, map[string]interface{}{
		"name":     user.Name,
		"password": user.Password,
		"token":    user.Token,
	}).Err()
}

func (r *RedisRepo) FetchUsers(ctx context.Context) (users []*domain.User, err error) {
	ids, err := r.client.SMembers(ctx, "users:").Result()
	if err == redis.Nil {
		return []*domain.User{}, nil
	}
	if err != nil {
		return
	}
	for _, id := range ids {
		id := strings.TrimPrefix(id, "user:")
		res := r.client.HGetAll(ctx, "user:"+id)
		if res.Err() == redis.Nil {
			continue
		}
		if res.Err() != nil {
			continue
		}
		var user domain.User
		if err := res.Scan(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	return
}

func (r *RedisRepo) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	id, err := r.client.HGet(ctx, "tokens:", token).Result()
	if err != nil {
		return &domain.User{}, err
	}
	id = strings.TrimPrefix(id, "user:")
	return r.FetchUserByID(ctx, id)
}

func (r *RedisRepo) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	res := r.client.HGetAll(ctx, "user:"+ID)
	if res.Err() != nil {
		return &domain.User{}, err
	}

	user = &domain.User{}
	if err := res.Scan(user); err != nil {
		return &domain.User{}, err
	}
	return
}
