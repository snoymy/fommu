package repo

import (
	"app/internal/api/core/entity"
	"context"
)

type SessionsRepo interface {
    FindSessionByID(ctx context.Context, id string) (*entity.SessionEntity, error)
    CreateSession(ctx context.Context, session *entity.SessionEntity) error
    UpdateSession(ctx context.Context, session *entity.SessionEntity) error
}
