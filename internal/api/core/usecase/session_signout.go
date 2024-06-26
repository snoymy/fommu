package usecase

import (
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"context"
)

type SignOutUsecase struct {
    sessionRepo repo.SessionsRepo
}

func NewSignOutUsecase(sessionRepo repo.SessionsRepo) *SignOutUsecase{
    return &SignOutUsecase{
        sessionRepo: sessionRepo,
    }
}

func (uc *SignOutUsecase) Exec(ctx context.Context, sessionId string) error {
    if sessionId == "" {
        return appstatus.BadValue("No session provided.")
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return err
    }

    if session == nil {
        return appstatus.BadValue("Session not found.")
    }

    if err := uc.sessionRepo.DeleteSession(ctx, sessionId); err != nil {
        return err
    }

    return nil
}
