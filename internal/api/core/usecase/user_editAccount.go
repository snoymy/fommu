package usecase

import (
	"app/internal/api/core/dto"
	"app/internal/api/core/repo"
	"app/internal/api/core/validator"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/log"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type EditAccountUsecase struct {
    userRepo repo.UsersRepo `injectable:""`
}

func NewEditAccountUsecase() *EditAccountUsecase {
    return &EditAccountUsecase{}
}

func (uc *EditAccountUsecase) Exec(ctx context.Context, username string, account dto.UserAccountDTO) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    hasUpdate := false

    if username == "" {
        log.Warn(ctx, "Cannot edit account, username is empty.")
        return appstatus.BadValue("username is empty.")
    }

    log.Info(ctx, "Check if user is exist")
    user, err := uc.userRepo.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }
    if user == nil {
        log.Info(ctx, "User not found")
        return appstatus.NotFound("user not found.")
    }

    if account.Email != nil {
        log.Info(ctx, "Update email")
        email := account.Email.ValueOrZero()
        if email != user.Email.ValueOrZero() {
            if account.CurrentPassword == nil {
                log.Warn(ctx, "Password is empty")
                return appstatus.InvalidCredential("Wrong password")
            }

            currentPassword := account.CurrentPassword.ValueOrZero()
            currentPasswordHash := uc.hashPassword(currentPassword)

            log.Info(ctx, "Check password")
            if currentPasswordHash != user.PasswordHash.ValueOrZero() {
                log.Warn(ctx, "Password not match")
                return appstatus.InvalidCredential("Wrong password")
            }

            log.Info(ctx, "Validate new email")
            if err := validator.ValidateEmail(email); err != nil {
                log.Info(ctx, "Email validation failed: " + err.Error())
                return appstatus.BadEmail(err.Error())
            }

            log.Info(ctx, "Check if email is used")
            existUser, err := uc.userRepo.FindUserByEmail(ctx, email, config.Fommu.Domain)
            if err != nil {
                log.Error(ctx, "Error: " + err.Error())
                return appstatus.InternalServerError("Something went wrong")
            }
            if existUser != nil {
                log.Warn(ctx, "Email is already used")
                return appstatus.BadEmail("Email already used.")
            }

            log.Info(ctx, "Set new email")
            user.Email.Set(email)
            hasUpdate = true
        } else {
            log.Info(ctx, "Email has no change")
        }
    }

    if account.NewPassword != nil {
        log.Info(ctx, "Update password")
        if account.CurrentPassword == nil {
            log.Warn(ctx, "Password is empty")
            return appstatus.InvalidCredential("Wrong password")
        }

        currentPassword := account.CurrentPassword.ValueOrZero()
        currentPasswordHash := uc.hashPassword(currentPassword)

        log.Info(ctx, "Check password")
        if currentPasswordHash != user.PasswordHash.ValueOrZero() {
            log.Warn(ctx, "Password not match")
            return appstatus.InvalidCredential("Wrong password")
        }

        log.Info(ctx, "Validate new password")
        newPassword := account.NewPassword.ValueOrZero()
        if err := validator.ValidatePassword(newPassword); err != nil {
            log.Info(ctx, "Password validation failed: " + err.Error())
            return appstatus.BadPassword(err.Error())
        }
        newPasswordHash := uc.hashPassword(newPassword)
        if newPasswordHash != currentPasswordHash {
            log.Info(ctx, "Set password")
            user.PasswordHash.Set(newPasswordHash)
            hasUpdate = true
        } else {
            log.Info(ctx, "Password has no change")
        }
    }

    if hasUpdate {
        user.UpdateAt.Set(time.Now().UTC())
        if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return appstatus.InternalServerError("Something went wrong")
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
