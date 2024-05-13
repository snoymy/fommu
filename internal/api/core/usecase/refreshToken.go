package usecase

import (
	"app/internal/appstatus"
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/utils"
	"context"
	"time"
)

type RefreshTokenUsecase struct {
    sessionRepo repo.SessionsRepo
}

func NewRefreshTokenUsecase(sessionRepo repo.SessionsRepo) *RefreshTokenUsecase {
    return &RefreshTokenUsecase{
        sessionRepo: sessionRepo,
    }
}

func (uc *RefreshTokenUsecase) Exec(ctx context.Context, sessionId string, refreshToken string) (*entity.SessionEntity, error) {
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    if session == nil {
        return nil, appstatus.InvalidSession("Session not found.")
    }

    if session.RefreshToken != refreshToken {
        return nil, appstatus.InvalidToken("Invalid token.")
    }

    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        return nil, appstatus.InvalidToken("Token expired.")
    }

    // create session id
    newAccessToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }
    newRefreshToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }

    session.AccessToken = newAccessToken
    session.AccessExpireAt = time.Now().UTC().Add(time.Minute * 15)
    session.RefreshToken = newRefreshToken
    session.RefreshExpireAt = time.Now().UTC().AddDate(0, 0, 30)
    session.LastRefresh.Set(time.Now().UTC())
    // write session to db
    if err := uc.sessionRepo.UpdateSession(ctx, session); err != nil {
        return nil, err
    }
    // return session 
    return session, nil
}

