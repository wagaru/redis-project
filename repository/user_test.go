package repository

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/wagaru/redis-project/domain"
)

func generateMockRedis() (*RedisRepo, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	repo := &RedisRepo{
		client: client,
	}
	return repo, mock
}

func TestStoreUser(t *testing.T) {
	repo, mock := generateMockRedis()

	mockUser := domain.User{
		ID:       "fakeID",
		Token:    "fakeToken",
		Password: "fakePassWord",
	}

	mock.ExpectSAdd("users:", mockUser.ID).SetVal(1)
	mock.ExpectHSet("tokens:", mockUser.Token, mockUser.ID).SetVal(1)
	mock.ExpectHSet("user:"+mockUser.ID, map[string]interface{}{
		"name":     mockUser.Name,
		"password": mockUser.Password,
		"token":    mockUser.Token,
	}).SetVal(1)

	err := repo.StoreUser(context.Background(), &mockUser)

	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestFetchUsers(t *testing.T) {
	repo, mock := generateMockRedis()

	mockUsers := []domain.User{
		{
			ID:    "fakeUserID1",
			Name:  "fakeName1",
			Token: "token1",
		},
		{
			ID:    "fakeUserID2",
			Name:  "fakeName2",
			Token: "token2",
		},
	}

	expectedSetMembers := make([]string, 0)
	for _, mockUser := range mockUsers {
		expectedSetMembers = append(expectedSetMembers, mockUser.ID)
	}
	mock.ExpectSMembers("users:").SetVal(expectedSetMembers)

	for _, mockUser := range mockUsers {
		mock.ExpectHGetAll("user:" + mockUser.ID).SetVal(map[string]string{
			"name":  mockUser.Name,
			"token": mockUser.Token,
		})
	}

	users, err := repo.FetchUsers(context.Background())

	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	if len(users) < 2 {
		t.Errorf("invalid")
	}

	if users[0].Token != mockUsers[0].Token {
		t.Errorf("invalid data")
	}
}

func TestFetchUserById(t *testing.T) {
	repo, mock := generateMockRedis()

	mockUser := domain.User{
		ID:    "fakeUserID",
		Name:  "fakeName",
		Token: "fakeToken",
	}

	mock.ExpectHGetAll("user:" + mockUser.ID).SetVal(map[string]string{
		"name":  mockUser.Name,
		"token": mockUser.Token,
	})

	user, err := repo.FetchUserByID(context.Background(), mockUser.ID)
	if err != nil {
		t.Error(err)
	}

	if user.Name != mockUser.Name {
		t.Errorf("return user not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
