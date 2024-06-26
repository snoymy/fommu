package usecase

import (
	"app/internal/api/core/entity"
	"app/internal/api/core/repo"
	"app/internal/api/core/validator"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/types"
	"app/internal/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type SigninUsecase struct {
    userRepo repo.UsersRepo
    sessionRepo repo.SessionsRepo
}

func NewSigninUsecase(userRepo repo.UsersRepo, sessionRepo repo.SessionsRepo) *SigninUsecase {
    return &SigninUsecase{
        userRepo: userRepo,
        sessionRepo: sessionRepo,
    }
}

func (uc *SigninUsecase) Exec(ctx context.Context, email string, password string, clientData types.JsonObject) (*entity.SessionEntity, error) {
    // check email
    if email == "" {
        return nil, appstatus.BadLogin("Invalid email or password")
    }

    var user *entity.UserEntity = nil
    if err := validator.ValidateEmail(email); err == nil {
        // find user by email
        user, err = uc.userRepo.FindUserByEmail(ctx, email, config.Fommu.Domain)
        if err != nil {
            return nil, appstatus.InternalServerError(err.Error())
        }
    } else {
        user, err = uc.userRepo.FindUserByUsername(ctx, email, config.Fommu.Domain)
        if err != nil {
            return nil, appstatus.InternalServerError(err.Error())
        }
    }

    // return error if user is null
    if user == nil {
        return nil, appstatus.BadLogin("Invalid email or password")
    }

    // hash password
    // set password_hash
    passwordHash := uc.hashPassword(password)

    // return error if password not match
    if user.PasswordHash.ValueOrZero() != passwordHash {
        return nil, appstatus.BadLogin("Invalid email or password")
    }

    // create session id
    accessToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }
    refreshToken, err := utils.GenerateRandomKey(45)
    if err != nil {
        return nil, appstatus.InternalServerError(err.Error())
    }

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
    if err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
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
