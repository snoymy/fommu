package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"context"
	"time"
)

type GetTokenUsecase struct {
    sessionRepo repo.SessionsRepo
}

func NewGetTokenUsecase(sessionRepo repo.SessionsRepo) *GetTokenUsecase {
    return &GetTokenUsecase{
        sessionRepo: sessionRepo,
    }
}

func (uc *GetTokenUsecase) Exec(ctx context.Context, sessionId string) (*entity.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if sessionId == "" {
        return nil, appstatus.BadValue("No session provided.")
    }

    log.Info(ctx, "Check if session is exist")
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }

    if session == nil {
        log.Info(ctx, "Session not found")
        return nil, appstatus.InvalidSession("Session not found.")
    }

    log.Info(ctx, "Check if refresh expired")
    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        log.Info(ctx, "Refresh token is expired")
        return nil, appstatus.InvalidToken("Session expired.")
    }

    return session, nil
}
