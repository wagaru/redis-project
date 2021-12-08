package usecase

import (
	"context"

	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/repository"
)

type postusecase struct {
	repo repository.PostRepo
}
type PostUsecase interface {
	domain.PostUsecase
}

func NewPostUsecase(repo repository.PostRepo) PostUsecase {
	return &postusecase{
		repo: repo,
	}
}

func (u *postusecase) StorePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	post.Author = "user:" + user.ID
	return u.repo.StorePost(ctx, post)
}

func (u *postusecase) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	return u.repo.FetchPosts(ctx, params)
}
