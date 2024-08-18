package repoimpls

import (
	"app/internal/adapter/commands"
	"app/internal/adapter/queries"
	"app/internal/core/entities"
	"context"
)

type SessionsRepoImpl struct {
    queries *queries.Query `injectable:""`
    commands *commands.Command `injectable:""`
}

func NewSessionRepoImpl() *SessionsRepoImpl {
    return &SessionsRepoImpl{}
}

func (r *SessionsRepoImpl) CreateSession(ctx context.Context, session *entities.SessionEntity) error {
    err := r.commands.CreateSession(ctx, session)
    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) UpdateSession(ctx context.Context, session *entities.SessionEntity) error {
    err := r.commands.UpdateSession(ctx, session)
    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) FindSessionByID(ctx context.Context, id string) (*entities.SessionEntity, error) {
    session, err := r.queries.FindSessionById(ctx, id)
    if err != nil {
        return nil, err
    }

    return session, nil
}


func (r *SessionsRepoImpl) DeleteSession(ctx context.Context, id string) error {
    err := r.commands.DeleteSession(ctx, id)
    if err != nil {
        return err
    }

    return nil
}
