package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"context"
	"time"
)


type AuthUsecase struct {
    sessionRepo repo.SessionsRepo `injectable:""`
}

func NewAuthUsecase() *AuthUsecase {
    return &AuthUsecase{}
}

func (uc *AuthUsecase) Exec(ctx context.Context, sessionId string, accessToken string) (*entity.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if sessionId == "" {
        log.Info(ctx, "Session ID is empty")
        return nil, appstatus.InvalidSession("Session not found.")
    }
    if accessToken == "" {
        log.Info(ctx, "Access token is empty")
        return nil, appstatus.InvalidToken("Token not found.")
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

    log.Info(ctx, "Check if refresh token is expired")
    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        log.Info(ctx, "Refresh token is expired")
        return nil, appstatus.InvalidToken("Session expired.")
    }

    log.Info(ctx, "Check if access token is match")
    if session.AccessToken != accessToken {
        log.Info(ctx, "Access token not match")
        return nil, appstatus.InvalidToken("Invalid token.")
    }

    log.Info(ctx, "Check if access token is expired")
    if session.AccessExpireAt.Compare(time.Now().UTC()) <= -1 {
        log.Info(ctx, "Access token is expired")
        return nil, appstatus.InvalidToken("Token expired.")
    }

    return session, nil
}
