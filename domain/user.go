package domain

import "context"

type User struct {
	ID       string `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
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
