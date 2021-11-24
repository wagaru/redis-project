package usecase

import (
	"context"
	"encoding/hex"
	"math/rand"

	"github.com/google/uuid"
	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/repository"
)

type usecase struct {
	repo repository.Repo
}

type Usecase interface {
	domain.PostUsecase
	domain.UserUsecase
}

func NewUsecase(repo repository.Repo) Usecase {
	return &usecase{
		repo: repo,
	}
}

// func (u *usecase) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Post, string, error) {
// 	return u.repo.Fetch(ctx, cursor, num)
// }

// func (u *usecase) GetById(ctx context.Context, id int64) (domain.Post, error) {
// 	return u.repo.GetById(ctx, id)
// }

func (u *usecase) StorePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	post.ID = uuid.NewString()
	post.Author = "user:" + user.ID
	return u.repo.StorePost(ctx, post)
}

func (u *usecase) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	return u.repo.FetchPosts(ctx, params)
}

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (u *usecase) StoreUser(ctx context.Context, user *domain.User) (err error) {
	user.ID = uuid.NewString()
	user.Token = GenerateToken()
	return u.repo.StoreUser(ctx, user)
}

func (u *usecase) FetchUsers(ctx context.Context) (user []*domain.User, err error) {
	return u.repo.FetchUsers(ctx)
}

func (u *usecase) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	return u.repo.FetchUserByToken(ctx, token)
}

func (u *usecase) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	return u.repo.FetchUserByID(ctx, ID)
}
