package repository_mock

import (
	"app/domain/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
    mock.Mock
}

func (m *UserRepoMock) FindUserByUsername(ctx context.Context, username string) (*entity.UserEntity, error) {
    args := m.Called(username)
    var user *entity.UserEntity
    if u, ok := args.Get(0).(*entity.UserEntity); ok {
        user = u
    } else {
        user = nil
    }
    
    return user, args.Error(1)
}

func (m *UserRepoMock) FindUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
    args := m.Called(email)
    var user *entity.UserEntity
    if u, ok := args.Get(0).(*entity.UserEntity); ok {
        user = u
    } else {
        user = nil
    }
    
    return user, args.Error(1)
}

func (m *UserRepoMock) FindResource(ctx context.Context, resource string, domain string) (*entity.UserEntity, error) {
    args := m.Called(resource, domain)
    var user *entity.UserEntity
    if u, ok := args.Get(0).(*entity.UserEntity); ok {
        user = u
    } else {
        user = nil
    }
    
    return user, args.Error(1)
}

func (m *UserRepoMock) CreateUser(ctx context.Context, user *entity.UserEntity) error {
    args := m.Called(user)
    return args.Error(1)
}
