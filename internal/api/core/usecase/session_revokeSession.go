package usecase

import (
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"context"
)

type RevokeSessionUsecase struct {
    sessionRepo repo.SessionsRepo
}

func NewRevokeSessionUsecase(sessionRepo repo.SessionsRepo) *RevokeSessionUsecase {
    return &RevokeSessionUsecase{
        sessionRepo: sessionRepo,
    }
}

func (uc *RevokeSessionUsecase) Exec(ctx context.Context, currentSessionId string, sessionId string) error {
    if currentSessionId == "" {
        return appstatus.BadValue("No session provided.")
    }

    if sessionId == "" {
        return appstatus.BadValue("No session to revoke.")
    }

    currentSession, err := uc.sessionRepo.FindSessionByID(ctx, currentSessionId)
    if err != nil {
        return err
    }

    if currentSession == nil {
        return appstatus.InvalidSession("Session not found.")
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return err
    }

    if session == nil {
        return appstatus.BadValue("Session not found.")
    }

    if currentSession.Owner != session.Owner {
        return appstatus.BadValue("You are not session owner.")
    }

    if err := uc.sessionRepo.DeleteSession(ctx, sessionId); err != nil {
        return err
    }

    return nil
}
