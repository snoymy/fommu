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
	"net/url"
	"time"

	"github.com/google/uuid"
)

type SignupUsecase struct {
    userRepository repos.UsersRepo `injectable:""`
}

func NewSignupUsecase() *SignupUsecase {
    return &SignupUsecase{}
}

func (uc *SignupUsecase) Exec(ctx context.Context, email string, username string, password string) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Validate email.")
    if err := uc.validateEmail(ctx, email); err != nil {
        return err
    }

    log.Info(ctx, "Validate username.")
    if err := uc.validateUsername(ctx, username); err != nil {
        return err
    }

    log.Info(ctx, "Validate password.")
    if err := uc.validatePassword(password); err != nil {
        return err
    }

    log.Info(ctx, "Create user entity.")
    user, err := uc.createUser(email, username, password)
    if err != nil {
        return err
    }

    log.Info(ctx, "Write user to database.")
    if err := uc.userRepository.CreateUser(ctx, user); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Somethig went wrong")
    }

    return nil
}

func (uc *SignupUsecase) createUser(email string, username string, password string) (*entities.UserEntity, error) {
    user := entities.NewUserEntity()
    user.ID = uuid.New().String()
    user.Email.Set(email)
    user.Username = username

    user.PasswordHash.Set(passwordutil.HashPassword(password))
    user.Displayname = username
    user.Remote = false
    user.Discoverable = true 
    user.AutoApproveFollower = false
    user.Domain = config.Fommu.Domain
    user.Preference.SetNull()
    user.Attachment.Set(types.JsonArray{})
    user.Tag.Set(types.JsonArray{})


    userUrl, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        return nil, err
    }

    inboxURL, err := url.JoinPath(userUrl, "inbox")
    if err != nil {
        return nil, err
    }

    outbox, err := url.JoinPath(userUrl, "outbox")
    if err != nil {
        return nil, err
    }

    followersURL, err := url.JoinPath(userUrl, "followers")
    if err != nil {
        return nil, err
    }

    followingURL, err := url.JoinPath(userUrl, "following")
    if err != nil {
        return nil, err
    }

    const bitSize = 4096
    privateKeyByte, publicKeyByte, err := keygenutil.GenerateKeyPair(bitSize)
    if err != nil {
        return nil, appstatus.InternalServerError("Somethig went wrong")
    }

    user.ActorId = userUrl
    user.URL = userUrl
    user.InboxURL = inboxURL
    user.OutboxURL = outbox
    user.FollowersURL = followersURL
    user.FollowingURL = followingURL

    user.PrivateKey.Set(string(privateKeyByte))
    user.PublicKey = string(publicKeyByte)
    user.JoinAt.Set(time.Now().UTC())
    user.CreateAt = time.Now().UTC()
    user.Status = entities.UserStatusActive

    return user, nil
}

func (uc *SignupUsecase) validateUsername(ctx context.Context, username string) error {
    if err := validator.ValidateUsername(username); err != nil {
        return appstatus.BadUsername(err.Error())
    }

    existUser, err := uc.userRepository.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        return appstatus.InternalServerError("Somethig went wrong")
    }
    if existUser != nil {
        return appstatus.BadUsername("Username already used.")
    }

    return nil
}

func (uc *SignupUsecase) validateEmail(ctx context.Context, email string) error {
    if err := validator.ValidateEmail(email); err != nil {
        return appstatus.BadEmail(err.Error())
    }

    existUser, err := uc.userRepository.FindUserByEmail(ctx, email, config.Fommu.Domain)
    if err != nil {
        return appstatus.InternalServerError("Somethig went wrong")
    }
    if existUser != nil {
        return appstatus.BadEmail("Email already used.")
    }

    return nil
}

func (uc *SignupUsecase) validatePassword(password string) error {
    if err := validator.ValidatePassword(password); err != nil {
        return appstatus.BadPassword(err.Error())
    }

    return nil
}
