package repoimpl

import (
	"app/internal/adapter/command"
	"app/internal/adapter/query"
	"app/internal/core/entities"
	"context"
)

type SessionsRepoImpl struct {
    query *query.Query `injectable:""`
    command *command.Command `injectable:""`
}

func NewSessionRepoImpl() *SessionsRepoImpl {
    return &SessionsRepoImpl{}
}

func (r *SessionsRepoImpl) CreateSession(ctx context.Context, session *entities.SessionEntity) error {
    err := r.command.CreateSession(ctx, session)
    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) UpdateSession(ctx context.Context, session *entities.SessionEntity) error {
    err := r.command.UpdateSession(ctx, session)
    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) FindSessionByID(ctx context.Context, id string) (*entities.SessionEntity, error) {
    session, err := r.query.FindSessionById(ctx, id)
    if err != nil {
        return nil, err
    }

    return session, nil
}


func (r *SessionsRepoImpl) DeleteSession(ctx context.Context, id string) error {
    err := r.command.DeleteSession(ctx, id)
    if err != nil {
        return err
    }

    return nil
}
