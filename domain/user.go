package domain

import "context"

type User struct {
	ID       string `json:"-" redis:"-"`
	Name     string `json:"name" redis:"name"`
	Password string `json:"password" redis:"password"`
	Token    string `json:"token" redis:"token"`
}

type UserUsecase interface {
	StoreUser(ctx context.Context, user *User) (err error)
	FetchUsers(ctx context.Context) (users []*User, err error)
	FetchUserByToken(ctx context.Context, token string) (user *User, err error)
	FetchUserByID(ctx context.Context, ID string) (user *User, err error)
}

type UserRepository interface {
	StoreUser(ctx context.Context, user *User) (err error)
	FetchUsers(ctx context.Context) (users []*User, err error)
	FetchUserByToken(ctx context.Context, token string) (user *User, err error)
	FetchUserByID(ctx context.Context, ID string) (user *User, err error)
}
