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

const VOTE_SCORE = 432 * 1000

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
	postTime, err := r.client.ZScore(ctx, "time:", "post:"+post.ID).Result()
	if err != nil {
		return err
	}
	if postTime < float64(time.Now().UnixMilli())-float64(time.Millisecond*1000*86400*7) {
		return errors.New("post cannot be voted.")
	}

	success, err := r.client.SAdd(ctx, "voted:"+post.ID, "user:"+user.ID).Result()
	if err != nil {
		return err
	}

	if success == 1 {
		err = r.client.ZIncrBy(ctx, "score:", float64(VOTE_SCORE), "post:"+post.ID).Err()
		if err != nil {
			return err
		}
		err = r.client.HIncrBy(ctx, "post:"+post.ID, "votes", 1).Err()
		if err != nil {
			return err
		}
	}
	return nil
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

func (r *RedisRepo) FetchGroupPosts(ctx context.Context, groupName string) ([]*domain.Post, error) {
	err := r.client.ZInterStore(ctx, "score:"+groupName, &redis.ZStore{
		Keys:      []string{"group:" + groupName, "score:"},
		Weights:   []float64{1, 1},
		Aggregate: "MAX",
	}).Err()
	if err != nil {
		return nil, err
	}
	err = r.client.ExpireAt(ctx, "score:"+groupName, time.Now().Add(time.Minute)).Err()
	if err != nil {
		return nil, err
	}
	return r.FetchPosts(ctx, &domain.PostQueryParams{
		QueryParams: domain.QueryParams{
			SortBy: [2]string{"score:" + groupName, "+"},
		},
	})
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
	case "score:programming":
		setKey = "score:programming"
	case "score:dancing":
		setKey = "score:dancing"
	default:
		setKey = "time:"
	}
	ids, err := r.client.ZRevRangeByScore(ctx, setKey, &redis.ZRangeBy{
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
		res := r.client.HGetAll(ctx, id)
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
	now := time.Now().UnixMilli()
	err := r.client.HSet(ctx, "post:"+post.ID, map[string]interface{}{
		"id":     post.ID,
		"title":  post.Title,
		"author": post.Author,
		"votes":  post.Votes,
		"time":   now,
	}).Err()
	if err != nil {
		return err
	}
	votedKey := "voted:" + post.ID
	err = r.client.SAdd(ctx, votedKey, post.Author).Err()
	if err != nil {
		return err
	}
	err = r.client.ExpireAt(ctx, votedKey, time.Now().Add(7*24*time.Hour)).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "time:", &redis.Z{
		Score:  float64(now),
		Member: "post:" + post.ID,
	}).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "score:", &redis.Z{
		Score:  float64(VOTE_SCORE) + float64(now),
		Member: "post:" + post.ID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo) GroupPost(ctx context.Context, post *domain.Post, groupName string) error {
	return r.client.SAdd(ctx, "group:"+groupName, "post:"+post.ID).Err()
}
