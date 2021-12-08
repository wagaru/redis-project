package usecase

import (
	"context"
	"errors"

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

func (u *postusecase) FetchPostByID(ctx context.Context, ID string) (*domain.Post, error) {
	return u.repo.FetchPostByID(ctx, ID)
}

func (u *postusecase) VotePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	if post == nil || user == nil {
		return errors.New("empty post or user is invalid")
	}
	return u.repo.VotePost(ctx, post, user)
}

func (u *postusecase) GroupPost(ctx context.Context, post *domain.Post, groupName string) error {
	if groupName == "" {
		return errors.New("invalid groupName")
	}
	return u.repo.GroupPost(ctx, post, groupName)
}

func (u *postusecase) FetchGroupPosts(ctx context.Context, groupName string) ([]*domain.Post, error) {
	return u.repo.FetchGroupPosts(ctx, groupName)
}
