package usecase

import (
	"app/internal/api/core/dto"
	"app/internal/api/core/repo"
	"app/internal/api/core/validator"
	"app/internal/appstatus"
	"app/internal/config"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type EditAccountUsecase struct {
    userRepo repo.UsersRepo
}

func NewEditAccountUsecase(userRepo repo.UsersRepo) *EditAccountUsecase {
    return &EditAccountUsecase{
        userRepo: userRepo,
    }
}

func (uc *EditAccountUsecase) Exec(ctx context.Context, username string, account dto.UserAccountDTO) error {
    hasUpdate := false

    if username == "" {
        return appstatus.BadValue("username is empty.")
    }

    user, err := uc.userRepo.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        return err
    }
    if user == nil {
        return appstatus.NotFound("user not found.")
    }

    if account.Email != nil {
        email := account.Email.ValueOrZero()
        if email != user.Email.ValueOrZero() {
            if account.CurrentPassword == nil {
                return appstatus.InvalidCredential("Wrong password")
            }

            currentPassword := account.CurrentPassword.ValueOrZero()
            currentPasswordHash := uc.hashPassword(currentPassword)

            if currentPasswordHash != user.PasswordHash.ValueOrZero() {
                return appstatus.InvalidCredential("Wrong password")
            }

            if err := validator.ValidateEmail(email); err != nil {
                return appstatus.BadEmail(err.Error())
            }

            existUser, err := uc.userRepo.FindUserByEmail(ctx, email, config.Fommu.Domain)
            if err != nil {
                return err
            }
            if existUser != nil {
                return appstatus.BadEmail("Email already used.")
            }

            user.Email.Set(email)
            hasUpdate = true
        }
    }

    if account.NewPassword != nil {
        if account.CurrentPassword == nil {
            return appstatus.InvalidCredential("Wrong password")
        }

        currentPassword := account.CurrentPassword.ValueOrZero()
        currentPasswordHash := uc.hashPassword(currentPassword)

        if currentPasswordHash != user.PasswordHash.ValueOrZero() {
            return appstatus.InvalidCredential("Wrong password")
        }

        if account.NewPassword != nil {
            newPassword := account.NewPassword.ValueOrZero()
            if err := validator.ValidatePassword(newPassword); err != nil {
                return appstatus.BadPassword(err.Error())
            }
            newPasswordHash := uc.hashPassword(newPassword)
            if newPasswordHash != currentPasswordHash {
                user.PasswordHash.Set(newPasswordHash)
                hasUpdate = true
            }
        }
    }

    if hasUpdate {
        user.UpdateAt.Set(time.Now().UTC())
        if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
            return err
        }
    }

    return nil
}

func (uc *EditAccountUsecase) hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    passwordHash := string(base64.StdEncoding.EncodeToString(h.Sum([]byte(password))))

    return passwordHash
}
