package usecase

import (
	"app/internal/application/fommu/repo"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
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

func (uc *GetIdentityUsecase) Exec(ctx context.Context, sessionId string) (*entities.UserEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    if sessionId == "" {
        log.Info(ctx, "Session id is empty")
        return nil, nil
    }

    log.Info(ctx, "Get session")
    session, err := uc.getSession(ctx, sessionId)
    if err != nil {
        return nil, err
    }

    log.Info(ctx, "Get session owner")
    user, err := uc.getSessionOwner(ctx, session.Owner)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (uc *GetIdentityUsecase) getSession(ctx context.Context, sessionId string) (*entities.SessionEntity, error) {
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

func (uc *GetIdentityUsecase) getSessionOwner(ctx context.Context, ownerId string) (*entities.UserEntity, error) {
    user, err := uc.userRepo.FindUserByID(ctx, ownerId)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }

    if user == nil {
        log.Info(ctx, "User not found")
        return nil, appstatus.NotFound("User not found")
    }

    return user, nil
}
