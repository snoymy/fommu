package usecase

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"app/internal/utils/keygenutil"
	"context"
	"time"
)

type RefreshTokenUsecase struct {
    sessionRepo repo.SessionsRepo `injectable:""`
}

func NewRefreshTokenUsecase() *RefreshTokenUsecase {
    return &RefreshTokenUsecase{}
}

func (uc *RefreshTokenUsecase) Exec(ctx context.Context, sessionId string, refreshToken string) (*entities.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Check if session is exist")
    session, err := uc.getSession(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    log.Info(ctx, "Check refresh token")
    if err := uc.checkRefreshToken(ctx, session, refreshToken); err != nil {
        return nil, err
    }

    log.Info(ctx, "Refresh session")
    if err := uc.refreshSession(ctx, session); err != nil {
        return nil, err
    }

    log.Info(ctx, "Write session to database")
    if err := uc.sessionRepo.UpdateSession(ctx, session); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }

    return session, nil
}

func (uc *RefreshTokenUsecase) getSession(ctx context.Context, sessionId string) (*entities.SessionEntity, error) {
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

func (uc *RefreshTokenUsecase) checkRefreshToken(ctx context.Context, session *entities.SessionEntity, refreshToken string) error {
    log.Info(ctx, "Check refresh token")
    if session.RefreshToken != refreshToken {
        log.Info(ctx, "Refresh token is not match")
        return appstatus.InvalidToken("Invalid token.")
    }

    log.Info(ctx, "Check if refresh expired")
    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        log.Info(ctx, "Refresh token is expired")
        return appstatus.InvalidToken("Token expired.")
    }
    return nil
}

func (uc *RefreshTokenUsecase) refreshSession(ctx context.Context, session *entities.SessionEntity) error {
    log.Info(ctx, "Create new access token")
    newAccessToken, err := keygenutil.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }
    log.Info(ctx, "Create new refresh token")
    newRefreshToken, err := keygenutil.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }

    log.Info(ctx, "Assign new token to entity")
    session.AccessToken = newAccessToken
    session.AccessExpireAt = time.Now().UTC().Add(time.Minute * 15)
    session.RefreshToken = newRefreshToken
    session.RefreshExpireAt = time.Now().UTC().AddDate(0, 0, 30)
    session.LastRefresh.Set(time.Now().UTC())

    return nil
}
