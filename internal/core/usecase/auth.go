package usecase

import (
	"app/internal/appstatus"
	"app/internal/core/repo"
	"context"
	"time"
)


type AuthUsecase struct {
    sessionRepo repo.SessionsRepo
}

func NewAuthUsecase(sessionRepo repo.SessionsRepo) *AuthUsecase {
    return &AuthUsecase{
        sessionRepo: sessionRepo,
    }
}

func (uc *AuthUsecase) Exec(ctx context.Context, sessionId string, accessToken string) error {
    if sessionId == "" {
        return appstatus.InvalidSession("Session not found.")
    }
    if accessToken == "" {
        return appstatus.InvalidToken("Token not found.")
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return err
    }

    if session == nil {
        return appstatus.InvalidSession("Session not found.")
    }

    if session.AccessToken != accessToken {
        return appstatus.InvalidToken("Invalid token.")
    }

    if session.AccessExpireAt.Compare(time.Now().UTC()) <= -1 {
        return appstatus.InvalidToken("Token expired.")
    }


    return nil
}
