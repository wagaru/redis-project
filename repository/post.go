package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wagaru/redis-project/domain"
)

func (r *redisRepo) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
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

func (r *redisRepo) StorePost(ctx context.Context, post *domain.Post) error {
	now := time.Now()
	nowms := float64(now.UnixNano() / int64(time.Millisecond))
	err := r.client.HSet(ctx, "post:"+post.ID, map[string]interface{}{
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
	err = r.client.ExpireAt(ctx, votedKey, now.Add(20*24*time.Hour)).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "time:", &redis.Z{
		Score:  nowms,
		Member: post.ID,
	}).Err()
	if err != nil {
		return err
	}
	err = r.client.ZAdd(ctx, "score:", &redis.Z{
		Score:  float64(432) + nowms,
		Member: post.ID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
