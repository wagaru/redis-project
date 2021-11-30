package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/wagaru/redis-project/domain"
)

func (r *RedisRepo) StoreUser(ctx context.Context, user *domain.User) error {
	err := r.client.SAdd(ctx, "users:", user.ID).Err()
	if err != nil {
		return err
	}
	err = r.client.HSet(ctx, "tokens:", user.Token, user.ID).Err()
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
		res := r.client.HGetAll(ctx, "user:"+id)
		if res.Err() != nil {
			continue
		}
		user := &domain.User{}
		if err := res.Scan(user); err != nil {
			continue
		}
		users = append(users, user)
	}
	return
}

func (r *RedisRepo) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	id, err := r.client.HGet(ctx, "tokens:", token).Result()
	if err != nil {
		return &domain.User{}, err
	}
	return r.FetchUserByID(ctx, id)
}

func (r *RedisRepo) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	res := r.client.HGetAll(ctx, "user:"+ID)
	if res.Err() != nil {
		return &domain.User{}, err
	}
	user2 := domain.User{}
	if err := res.Scan(&user2); err != nil {
		return &domain.User{}, err
	}
	user = &user2
	return user, nil
}
