package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/wagaru/redis-project/domain"
)

func (u *usecase) StorePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	post.ID = uuid.NewString()
	post.Author = "user:" + user.ID
	return u.repo.StorePost(ctx, post)
}

func (u *usecase) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	return u.repo.FetchPosts(ctx, params)
}
