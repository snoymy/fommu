package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"context"
)

type GetIdentityUsecase struct {
    sessionRepo repo.SessionsRepo
    userRepo repo.UsersRepo
}

func NewGetIdentityUsecase(sessionRepo repo.SessionsRepo, userRepo repo.UsersRepo) *GetIdentityUsecase {
    return &GetIdentityUsecase{
        sessionRepo: sessionRepo,
        userRepo: userRepo,
    }
}

func (uc *GetIdentityUsecase) Exec(ctx context.Context, sessionId string) (*entity.UserEntity, error) {
    if sessionId == "" {
        return nil, nil
    }

    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    if session == nil {
        return nil, appstatus.InvalidSession("Session not found.")
    }

    user, err := uc.userRepo.FindUserByID(ctx, session.Owner)
    if err != nil {
        return nil, err
    }

    if user == nil {
        return nil, nil
    }

    return user, nil
}

