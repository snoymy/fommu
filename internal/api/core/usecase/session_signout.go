package usecase

import (
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
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
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if sessionId == "" {
        log.Info(ctx, "Session is empty")
        return appstatus.BadValue("No session provided.")
    }

    log.Info(ctx, "Check if session exist")
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return err
    }

    if session == nil {
        log.Info(ctx, "Session not found")
        return appstatus.BadValue("Session not found.")
    }

    log.Info(ctx, "Delete session")
    if err := uc.sessionRepo.DeleteSession(ctx, sessionId); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return err
    }

    return nil
}
