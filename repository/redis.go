package repository

import (
	"context"
	"time"

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

func (r *redisRepo) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	if len(params.SortBy) == 0 {
		params.SortBy[0] = "time"
		params.SortBy[1] = "+"
	}
	posts := make([]*domain.Post, 0)
	setKey := ""
	switch params.SortBy[0] {
	case "time":
		setKey = "time:"
	case "score":
		setKey = "score:"
	default:
		setKey = "time"
	}
	ids, err := r.client.ZRangeByScore(ctx, setKey, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(params.Page) * int64(params.PerPage),
		Count:  int64(params.PerPage),
	}).Result()
	if err != nil {
		return []*domain.Post{}, err
	}
	for _, id := range ids {
		res := r.client.HGetAll(ctx, "post:"+id)
		if res.Err() != nil {
			continue
		}
		var post domain.Post
		if err = res.Scan(&post); err != nil {
			continue
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *redisRepo) StorePost(ctx context.Context, post *domain.Post) error {
	now := time.Now()
	now_withms := float64(now.UnixNano()/int64(time.Millisecond)) / 1000
	err := r.client.HSet(ctx, "post:"+post.ID, map[string]interface{}{
		"title":  post.Title,
		"author": post.Author,
		"votes":  post.Votes,
		"time":   now_withms,
	}).Err()
	if err != nil {
		return err
	}
	votedKey := "voted:" + post.ID
	err = r.client.SAdd(ctx, votedKey, post.Author).Err()
	if err != nil {
		return err
	}
	err = r.client.ExpireAt(ctx, votedKey, now.Add(20*24*time.Hour)).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "time:", &redis.Z{
		Score:  now_withms,
		Member: post.ID,
	}).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "score:", &redis.Z{
		Score:  float64(432) + now_withms,
		Member: post.ID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

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
