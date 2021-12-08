package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"math/rand"

	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/repository"
)

type userusecase struct {
	repo repository.UserRepo
}

type UserUsecase interface {
	domain.UserUsecase
}

func NewUserUsecase(repo repository.UserRepo) UserUsecase {
	return &userusecase{
		repo: repo,
	}
}

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (u *userusecase) StoreUser(ctx context.Context, user *domain.User) (err error) {
	if user.Token == "" {
		user.Token = GenerateToken()
	}
	return u.repo.StoreUser(ctx, user)
}

func (u *userusecase) FetchUsers(ctx context.Context) (user []*domain.User, err error) {
	return u.repo.FetchUsers(ctx)
}

func (u *userusecase) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	if token == "" {
		return nil, errors.New("Empty token invalid")
	}
	return u.repo.FetchUserByToken(ctx, token)
}

func (u *userusecase) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	if ID == "" {
		return nil, errors.New("Empty ID invalid")
	}
	return u.repo.FetchUserByID(ctx, ID)
}
