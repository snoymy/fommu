package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
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
    if sessionId == "" {
        return nil, appstatus.BadValue("No session provided.")
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    if session == nil {
        return nil, appstatus.InvalidSession("Session not found.")
    }

    if session.RefreshExpireAt.Compare(time.Now().UTC()) <= -1 {
        return nil, appstatus.InvalidToken("Session expired.")
    }

    return session, nil
}
