// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/wagaru/redis-project/domain"
)

// PostUsecase is an autogenerated mock type for the PostUsecase type
type PostUsecase struct {
	mock.Mock
}

// FetchPosts provides a mock function with given fields: ctx, params
func (_m *PostUsecase) FetchPosts(ctx context.Context, params *domain.PostQueryParams) ([]*domain.Post, error) {
	ret := _m.Called(ctx, params)

	var r0 []*domain.Post
	if rf, ok := ret.Get(0).(func(context.Context, *domain.PostQueryParams) []*domain.Post); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.PostQueryParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorePost provides a mock function with given fields: ctx, post, user
func (_m *PostUsecase) StorePost(ctx context.Context, post *domain.Post, user *domain.User) error {
	ret := _m.Called(ctx, post, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Post, *domain.User) error); ok {
		r0 = rf(ctx, post, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}