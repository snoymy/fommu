package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/api/core/validator"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/log"
	"app/internal/types"
	"app/internal/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type SigninUsecase struct {
    userRepo    repo.UsersRepo    `injectable:""`
    sessionRepo repo.SessionsRepo `injectable:""`
}

func NewSigninUsecase() *SigninUsecase {
    return &SigninUsecase{}
}

func (uc *SigninUsecase) Exec(ctx context.Context, email string, password string, clientData types.JsonObject) (*entity.SessionEntity, error) {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    // check email
    log.Info(ctx, "Check user login")
    if email == "" {
        log.Info(ctx, "user login is empty")
        return nil, appstatus.BadLogin("Invalid email/username or password")
    }

    var user *entity.UserEntity = nil
    log.Info(ctx, "Check user login type")
    if err := validator.ValidateEmail(email); err == nil {
        // find user by email
        log.Info(ctx, "Login by email")
        user, err = uc.userRepo.FindUserByEmail(ctx, email, config.Fommu.Domain)
        if err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return nil, appstatus.InternalServerError(err.Error())
        }
    } else {
        log.Info(ctx, "Login by username")
        user, err = uc.userRepo.FindUserByUsername(ctx, email, config.Fommu.Domain)
        if err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return nil, appstatus.InternalServerError(err.Error())
        }
    }

    // return error if user is null
    if user == nil {
        log.Info(ctx, "User not found")
        return nil, appstatus.BadLogin("Invalid email/username or password")
    }

    log.Info(ctx, "Check password")
    // hash password
    // set password_hash
    passwordHash := uc.hashPassword(password)

    // return error if password not match
    if user.PasswordHash.ValueOrZero() != passwordHash {
        log.Info(ctx, "Invalid password")
        return nil, appstatus.BadLogin("Invalid email or password")
    }

    // create session id
    log.Info(ctx, "Create token")
    accessToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError(err.Error())
    }
    refreshToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError(err.Error())
    }

    log.Info(ctx, "Create session")
    session := entity.NewSessionEntity()
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
    // write session to db
    log.Info(ctx, "Write session to database")
    if err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, err
    }
    // return session 
    return session, nil
}

func (uc *SigninUsecase) hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    passwordHash := string(base64.StdEncoding.EncodeToString(h.Sum([]byte(password))))

    return passwordHash
}
