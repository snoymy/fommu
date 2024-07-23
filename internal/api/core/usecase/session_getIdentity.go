package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/appstatus"
	"app/internal/log"
	"context"
)

type GetIdentityUsecase struct {
    sessionRepo repo.SessionsRepo `injectable:""`
    userRepo    repo.UsersRepo    `injectable:""`
}

func NewGetIdentityUsecase() *GetIdentityUsecase {
    return &GetIdentityUsecase{}
}

func (uc *GetIdentityUsecase) Exec(ctx context.Context, sessionId string) (*entity.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if sessionId == "" {
        log.Info(ctx, "Session id is empty")
        return nil, nil
    }

    log.Info(ctx, "Check if session is exist")
    session, err := uc.sessionRepo.FindSessionByID(ctx, sessionId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    if session == nil {
        log.Info(ctx, "Session not found")
        return nil, appstatus.InvalidSession("Session not found.")
    }

    log.Info(ctx, "Get session owner")
    user, err := uc.userRepo.FindUserByID(ctx, session.Owner)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    if user == nil {
        log.Info(ctx, "User not found")
        return nil, nil
    }

    return user, nil
}

