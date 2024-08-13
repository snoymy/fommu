package usecase

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
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
    currentSession, err := uc.getSession(ctx, currentSessionId)
    if err != nil {
        return err
    }

    log.Info(ctx, "Check if target session is exist")
    session, err := uc.getSession(ctx, sessionId)
    if err != nil {
        return err
    }

    log.Info(ctx, "Check session owner")
    if uc.isSessionOwner(currentSession, session) {
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

func (uc *RevokeSessionUsecase) getSession(ctx context.Context, sessionId string) (*entities.SessionEntity, error) {
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    if session == nil {
        log.Info(ctx, "Session not found")
        return nil, appstatus.InvalidSession("Session not found.")
    }

    return session, nil
}

func (uc *RevokeSessionUsecase) isSessionOwner(ownerSession *entities.SessionEntity, targetSession *entities.SessionEntity) bool {
    if ownerSession.Owner == targetSession.Owner {
        return true
    }

    return false
}
