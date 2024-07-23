package usecase

import (
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"context"
)

type RevokeSessionUsecase struct {
    sessionRepo repo.SessionsRepo `injectable:""`
}

func NewRevokeSessionUsecase() *RevokeSessionUsecase {
    return &RevokeSessionUsecase{}
}

func (uc *RevokeSessionUsecase) Exec(ctx context.Context, currentSessionId string, sessionId string) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if currentSessionId == "" {
        log.Info(ctx, "Session id is empty")
        return appstatus.BadValue("No session provided.")
    }

    if sessionId == "" {
        log.Info(ctx, "Target session id is empty")
        return appstatus.BadValue("No session to revoke.")
    }

    log.Info(ctx, "Check if session is exist")
    currentSession, err := uc.sessionRepo.FindSessionByID(ctx, currentSessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }

    if currentSession == nil {
        log.Info(ctx, "Session not found")
        return appstatus.InvalidSession("Session not found.")
    }

    log.Info(ctx, "Check if target session is exist")
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }

    if session == nil {
        log.Info(ctx, "Target session not found")
        return appstatus.BadValue("Session not found.")
    }

    log.Info(ctx, "Check session owner")
    if currentSession.Owner != session.Owner {
        log.Info(ctx, "target session not found")
        return appstatus.BadValue("You are not session owner.")
    }

    log.Info(ctx, "Delete session")
    if err := uc.sessionRepo.DeleteSession(ctx, sessionId); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return err
    }

    return nil
}
