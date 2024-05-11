package usecase

import (
	"app/internal/appstatus"
	"app/internal/core/entity"
	"app/internal/core/repo"
	"app/internal/core/validator"
	"app/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

type SignupUsecase struct {
    userRepository repo.UsersRepo
}

func NewSignupUsecase(userRepository repo.UsersRepo) *SignupUsecase {
    return &SignupUsecase{
        userRepository: userRepository,
    }
}

func (uc *SignupUsecase) Exec(ctx context.Context, email string, username string, password string) error {
    var (
        existUser *entity.UserEntity
        err error
    )

    // validate username
    if err := validator.ValidateUsername(username); err != nil {
        return appstatus.BadUsername(err.Error())
    }

    // check if username is used
    existUser, err = uc.userRepository.FindUserByUsername(ctx, username)
    if err != nil {
        return appstatus.InternalServerError(err.Error())
    }
    if existUser != nil {
        return appstatus.BadUsername("Username already used.")
    }

    // validate email
    if err := validator.ValidateEmail(email); err != nil {
        return appstatus.BadEmail(err.Error())
    }

    // check if email is used
    existUser, err = uc.userRepository.FindUserByEmail(ctx, email)
    if err != nil {
        return err
    }
    if existUser != nil {
        return appstatus.BadEmail("Email already used.")
    }

    // validate password
    if err := validator.ValidatePassword(password); err != nil {
        return appstatus.BadPassword(err.Error())
    }
    
    user := entity.NewUserEntity()
    // set id
    user.ID = uuid.New().String()
    // set email
    user.Email = email
    // set username
    user.Username = username

    // hash password
    // set password_hash
    user.PasswordHash = utils.HashPassword(password)

    // set display name
    user.Displayname = username
    // set url
    // user.URL, err = url.JoinPath(config.Fommu.URL, "users", user.Username)
    // if err != nil {
    //     return appstatus.InternalServerError(err.Error())
    // }
    // set remote
    user.Remote = false
    // set discoverable
    user.Discoverable = true 
    // set auto_approve_follower
    user.AutoApproveFollower = false

    // generate key
    const bitSize = 4096
    privateKeyByte, publicKeyByte, err := utils.GenerateKeyPair(bitSize)
    if err != nil {
        return appstatus.InternalServerError(err.Error())
    }

    // set private key
    user.PrivateKey = string(privateKeyByte)
    // set public key
    user.PublicKey = string(publicKeyByte)

    // set create date
    user.CreateAt = time.Now().UTC()
    // set active
    user.Status = "active"

    // write user to db
    if err := uc.userRepository.CreateUser(ctx, user); err != nil {
        return appstatus.InternalServerError(err.Error())
    }

    return nil
}
