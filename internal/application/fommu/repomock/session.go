package repository_mock

import (
	"app/internal/core/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

type SessionRepoMock struct {
    mock.Mock
}

func (m *SessionRepoMock) CreateSession(ctx context.Context, session *entity.SessionEntity) error {
    args := m.Called(session)
    return args.Error(1)
}

