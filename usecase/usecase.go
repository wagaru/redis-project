package usecase

import (
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
