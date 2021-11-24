package repository

import (
	"context"

	"github.com/wagaru/redis-project/domain"
)

func (r *redisRepo) StoreUser(ctx context.Context, user *domain.User) error {
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

func (r *redisRepo) FetchUsers(ctx context.Context) (users []*domain.User, err error) {
	ids, err := r.client.SMembers(ctx, "users:").Result()
	if err != nil {
		return
	}
	for _, id := range ids {
		user, err := r.client.HGetAll(ctx, "user:"+id).Result()
		if err != nil {
			continue
		}
		users = append(users, &domain.User{
			Name:  user["name"],
			Token: user["token"],
		})
	}
	return
}

func (r *redisRepo) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	id, err := r.client.HGet(ctx, "tokens:", token).Result()
	if err != nil {
		return &domain.User{}, err
	}
	return r.FetchUserByID(ctx, id)
}

func (r *redisRepo) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	data, err := r.client.HGetAll(ctx, "user:"+ID).Result()
	if err != nil {
		return &domain.User{}, err
	}
	user = &domain.User{
		ID:    ID,
		Name:  data["name"],
		Token: data["token"],
	}
	return user, nil
}
