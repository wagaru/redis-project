package domain

import (
	"context"
	"time"
)

type Post struct {
	ID     string    `json:"-" redis:"id"`
	Title  string    `json:"title" redis:"title"`
	Author string    `json:"author" redis:"author"`
	Votes  int       `json:"votes" redis:"votes"`
	Time   time.Time `json:"time"`
}

type PostQueryParams struct {
	QueryParams
}

type PostUsecase interface {
	// Fetch(ctx context.Context, cursor string, num int64) ([]Post, string, error)
	// GetById(ctx context.Context, id int64) (Post, error)
	StorePost(ctx context.Context, post *Post, user *User) error
	FetchPosts(ctx context.Context, params *PostQueryParams) ([]*Post, error)
	FetchPostByID(ctx context.Context, ID string) (*Post, error)
	VotePost(ctx context.Context, post *Post, user *User) error
}

type PostRepository interface {
	// Fetch(ctx context.Context, cursor string, num int64) ([]Post, string, error)
	// GetById(ctx context.Context, id int64) (Post, error)
	StorePost(ctx context.Context, post *Post) error
	FetchPosts(ctx context.Context, params *PostQueryParams) ([]*Post, error)
	FetchPostByID(ctx context.Context, ID string) (*Post, error)
	VotePost(ctx context.Context, post *Post, user *User) error
}
