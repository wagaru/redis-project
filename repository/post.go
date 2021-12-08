package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wagaru/redis-project/domain"
)

type PostRepo interface {
	domain.PostRepository
}

func NewRedisPostRepo(addr string) (PostRepo, error) {
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

func (r *RedisRepo) VotePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	//檢查是否已過期，過期不給投

	//更新 score:

	//更新 post:xx 上的　votes
}

func (r *RedisRepo) FetchPostByID(ctx context.Context, ID string) (*domain.Post, error) {
	res := r.client.HGetAll(ctx, "post:"+ID)
	if res.Err() == redis.Nil {
		return nil, errors.New("no matching post")
	}
	var post domain.Post
	if err := res.Scan(&post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *RedisRepo) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	if len(params.SortBy) == 0 {
		params.SortBy[0] = "time"
		params.SortBy[1] = "+"
	}
	if params.PerPage == 0 {
		params.PerPage = 10
	}
	posts := make([]*domain.Post, 0)
	setKey := ""
	switch params.SortBy[0] {
	case "time":
		setKey = "time:"
	case "score":
		setKey = "score:"
	default:
		setKey = "time:"
	}
	fmt.Println(setKey)
	ids, err := r.client.ZRangeByScore(ctx, setKey, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(params.Page) * int64(params.PerPage),
		Count:  int64(params.PerPage),
	}).Result()
	if err != nil {
		fmt.Println(err)
		return []*domain.Post{}, err
	}
	for _, id := range ids {
		res := r.client.HGetAll(ctx, "post:"+id)
		if res.Err() != nil {
			fmt.Println(err)
			continue
		}
		var post domain.Post
		if err = res.Scan(&post); err != nil {
			fmt.Println(err)
			continue
		}
		result, _ := res.Result()
		timems, _ := strconv.Atoi(result["time"])
		post.Time = time.Unix(int64(timems/1000), int64(timems%1000))
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *RedisRepo) StorePost(ctx context.Context, post *domain.Post) error {
	if post.ID == "" {
		res, err := r.client.Incr(ctx, "post:").Result()
		if err != nil {
			return err
		}
		post.ID = strconv.FormatInt(res, 10)
	}
	now := time.Now()
	nowms := float64(now.UnixNano() / int64(time.Millisecond))
	err := r.client.HSet(ctx, "post:"+post.ID, map[string]interface{}{
		"id":     post.ID,
		"title":  post.Title,
		"author": post.Author,
		"votes":  post.Votes,
		"time":   nowms,
	}).Err()
	if err != nil {
		return err
	}
	votedKey := "voted:" + post.ID
	err = r.client.SAdd(ctx, votedKey, post.Author).Err()
	if err != nil {
		return err
	}
	err = r.client.ExpireAt(ctx, votedKey, now.Add(7*24*time.Hour)).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "time:", &redis.Z{
		Score:  nowms,
		Member: "post:" + post.ID,
	}).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "score:", &redis.Z{
		Score:  float64(432) + nowms,
		Member: "post:" + post.ID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
