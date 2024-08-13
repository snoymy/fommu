package ports

import (
	"app/internal/core/entities"
	"context"
)

type SessionsRepo interface {
    FindSessionByID(ctx context.Context, id string) (*entities.SessionEntity, error)
    CreateSession(ctx context.Context, session *entities.SessionEntity) error
    UpdateSession(ctx context.Context, session *entities.SessionEntity) error
    DeleteSession(ctx context.Context, id string) error
}
