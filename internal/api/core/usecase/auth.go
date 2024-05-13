package usecase

import (
	"app/internal/appstatus"
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
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

func (uc *AuthUsecase) Exec(ctx context.Context, sessionId string, accessToken string) (*entity.SessionEntity, error) {
    if sessionId == "" {
        return nil, appstatus.InvalidSession("Session not found.")
    }
    if accessToken == "" {
        return nil, appstatus.InvalidToken("Token not found.")
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    if session == nil {
        return nil, appstatus.InvalidSession("Session not found.")
    }

    if session.AccessToken != accessToken {
        return nil, appstatus.InvalidToken("Invalid token.")
    }

    if session.AccessExpireAt.Compare(time.Now().UTC()) <= -1 {
        return nil, appstatus.InvalidToken("Token expired.")
    }


    return session, nil
}
