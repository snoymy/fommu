package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"app/internal/utils"
	"context"
	"time"
)

type RefreshTokenUsecase struct {
    sessionRepo repo.SessionsRepo `injectable:""`
}

func NewRefreshTokenUsecase() *RefreshTokenUsecase {
    return &RefreshTokenUsecase{}
}

func (uc *RefreshTokenUsecase) Exec(ctx context.Context, sessionId string, refreshToken string) (*entity.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

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

    log.Info(ctx, "Check refresh token")
    if session.RefreshToken != refreshToken {
        log.Info(ctx, "Refresh token is not match")
        return nil, appstatus.InvalidToken("Invalid token.")
    }

    log.Info(ctx, "Check if refresh expired")
    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        log.Info(ctx, "Refresh token is expired")
        return nil, appstatus.InvalidToken("Token expired.")
    }

    // create session id
    log.Info(ctx, "Create new access token")
    newAccessToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }
    log.Info(ctx, "Create new refresh token")
    newRefreshToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    log.Info(ctx, "Assign new token to entity")
    session.AccessToken = newAccessToken
    session.AccessExpireAt = time.Now().UTC().Add(time.Minute * 15)
    session.RefreshToken = newRefreshToken
    session.RefreshExpireAt = time.Now().UTC().AddDate(0, 0, 30)
    session.LastRefresh.Set(time.Now().UTC())
    // write session to db
    log.Info(ctx, "Write session to database")
    if err := uc.sessionRepo.UpdateSession(ctx, session); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }
    // return session 
    return session, nil
}

