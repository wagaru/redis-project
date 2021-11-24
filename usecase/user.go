package usecase

import (
	"context"
	"encoding/hex"
	"math/rand"

	"github.com/google/uuid"
	"github.com/wagaru/redis-project/domain"
)

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
