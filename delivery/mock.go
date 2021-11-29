package delivery

import (
	"context"

	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/usecase"
)

type MockUsecase struct {
}

func NewMockUsecase() usecase.Usecase {
	return &MockUsecase{}
}

func (usecase *MockUsecase) StorePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	return nil
}

func (usecase *MockUsecase) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	return []*domain.Post{}, nil
}

func (usecase *MockUsecase) StoreUser(ctx context.Context, user *domain.User) (err error) {
	return nil
}

func (usecase *MockUsecase) FetchUsers(ctx context.Context) (users []*domain.User, err error) {
	return []*domain.User{}, nil
}

func (usecase *MockUsecase) FetchUserByToken(ctx context.Context, token string) (user *domain.User, err error) {
	return &domain.User{}, nil
}

func (usecase *MockUsecase) FetchUserByID(ctx context.Context, ID string) (user *domain.User, err error) {
	return &domain.User{}, nil
}
