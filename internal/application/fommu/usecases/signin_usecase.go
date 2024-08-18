package usecases

import (
	"app/internal/application/fommu/repos"
	"app/internal/application/fommu/validator"
	"app/internal/config"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/core/types"
	"app/internal/log"
	"app/internal/utils/keygenutil"
	"app/internal/utils/passwordutil"
	"context"
	"time"

	"github.com/google/uuid"
)

type SigninUsecase struct {
    userRepo    repos.UsersRepo    `injectable:""`
    sessionRepo repos.SessionsRepo `injectable:""`
}

func NewSigninUsecase() *SigninUsecase {
    return &SigninUsecase{}
}

func (uc *SigninUsecase) Exec(ctx context.Context, email string, password string, clientData types.JsonObject) (*entities.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Check user login")
    if email == "" {
        log.Info(ctx, "user login is empty")
        return nil, appstatus.BadLogin("Invalid email/username or password")
    }
    
    user, err := uc.getUserLogin(ctx, email)
    if err != nil {
        return nil, err
    }

    log.Info(ctx, "Check password")
    if !uc.isPasswordMatch(user, password) {
        return nil, appstatus.BadLogin("Invalid email or password")
    }

    log.Info(ctx, "Create session")
    session, err := uc.createSession(user, clientData)
    if err != nil {
        return nil, err
    }

    log.Info(ctx, "Write session to database")
    if err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }

    return session, nil
}

func (uc *SigninUsecase) createSession(user *entities.UserEntity, clientData types.JsonObject) (*entities.SessionEntity, error) {
    accessToken, err := keygenutil.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }
    refreshToken, err := keygenutil.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }

    session := entities.NewSessionEntity()
    session.ID = uuid.New().String()
    session.AccessToken = accessToken
    session.AccessExpireAt = time.Now().UTC().Add(time.Minute * 15)
    session.RefreshToken = refreshToken
    session.RefreshExpireAt = time.Now().UTC().AddDate(0, 0, 30)
    session.Owner = user.ID
    if clientData == nil {
        session.Metadata.SetNull()
    } else {
        session.Metadata.Set(clientData)
    }
    session.LoginAt = time.Now().UTC()
    session.LastRefresh.SetNull()

    return session, nil
}

func (uc *SigninUsecase) getUserLogin(ctx context.Context, email string) (*entities.UserEntity, error) {
    var user *entities.UserEntity
    if err := validator.ValidateEmail(email); err == nil {
        var err error
        user, err = uc.userRepo.FindUserByEmail(ctx, email, config.Fommu.Domain)
        if err != nil {
            return nil, appstatus.InternalServerError(err.Error())
        }
    } else {
        var err error
        user, err = uc.userRepo.FindUserByUsername(ctx, email, config.Fommu.Domain)
        if err != nil {
            return nil, appstatus.InternalServerError(err.Error())
        }
    }

    if user == nil {
        return nil, appstatus.BadLogin("Invalid email/username or password")
    }

    return user, nil
}

func (uc *SigninUsecase) isPasswordMatch(user *entities.UserEntity, password string) bool {
    passwordHash := passwordutil.HashPassword(password)
    if user.PasswordHash.ValueOrZero() == passwordHash {
        return true
    }
    return false
}
