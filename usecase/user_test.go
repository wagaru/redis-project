package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wagaru/redis-project/domain"
	"github.com/wagaru/redis-project/domain/mocks"
)

func TestStoreUser(t *testing.T) {
	repo := new(mocks.UserRepository)
	usecase := NewUserUsecase(repo)
	user := &domain.User{}
	repo.On("StoreUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool { return user.Token != "" })).Return(nil)
	err := usecase.StoreUser(context.Background(), user)
	if err != nil {
		t.Errorf("failed")
	}
}

func TestFetchUserByToken(t *testing.T) {
	repo := new(mocks.UserRepository)
	usecase := NewUserUsecase(repo)

	t.Run("Error", func(t *testing.T) {
		_, err := usecase.FetchUserByToken(context.Background(), "")
		if err == nil {
			t.Errorf("error")
		}
	})
	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			Token: "fakeToken",
		}
		repo.On("FetchUserByToken", mock.Anything, user.Token).Return(user, nil)
		res, err := usecase.FetchUserByToken(context.Background(), user.Token)
		if err != nil {
			t.Error("error")
		}
		if res.Token != user.Token {
			t.Error("invalid")
		}
	})
}

func TestFetchUserByID(t *testing.T) {
	repo := new(mocks.UserRepository)
	usecase := NewUserUsecase(repo)

	t.Run("Error", func(t *testing.T) {
		_, err := usecase.FetchUserByID(context.Background(), "")
		if err == nil {
			t.Errorf("error")
		}
	})
	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID: "fakeID",
		}
		repo.On("FetchUserByID", mock.Anything, user.ID).Return(user, nil)
		res, err := usecase.FetchUserByID(context.Background(), user.ID)
		if err != nil {
			t.Error("error")
		}
		if res.ID != user.ID {
			t.Error("invalid")
		}
	})
}
