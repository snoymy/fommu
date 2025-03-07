package usecases

import (
	"app/internal/application/fommu/dto"
	"app/internal/application/fommu/repos"
	"app/internal/application/fommu/validator"
	"app/internal/config"
	"app/internal/application/appstatus"
	"app/internal/core/entities"
	"app/internal/log"
	"app/internal/utils/passwordutil"
	"context"
	"errors"
	"time"
)

type EditAccountUsecase struct {
    userRepo repos.UsersRepo `injectable:""`
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

    log.Info(ctx, "Get user")
    user, err := uc.getUser(ctx, username)
    if err != nil {
        return err
    }

    if account.Email != nil {
        log.Info(ctx, "Update email")
        email := account.Email.ValueOrZero()
        if err := uc.updateEmail(user, email); err != nil {
            return err
        }
        hasUpdate = true
    }

    if account.NewPassword != nil {
        log.Info(ctx, "Update password")
        if err := uc.updatePassword(user, account.NewPassword.ValueOrZero()); err != nil {
            return err
        }
        hasUpdate = true
    }

    if account.Discoverable != nil {
        log.Info(ctx, "Update discoverable")
        if err := uc.updateDiscoverable(user, account.Discoverable.ValueOrZero()); err != nil {
            return err
        }
        hasUpdate = true
    }

    if hasUpdate {
        currentPassword := account.CurrentPassword.ValueOrZero()
        if err := uc.checkPassword(ctx, user, currentPassword); err != nil {
            return err
        }
        user.UpdateAt.Set(time.Now().UTC())
        if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return appstatus.InternalServerError("Something went wrong")
        }
    }

    return nil
}

func (uc *EditAccountUsecase) getUser(ctx context.Context, username string) (*entities.UserEntity, error) {
    user, err := uc.userRepo.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return nil, appstatus.InternalServerError("Something went wrong")
    }
    if user == nil {
        log.Info(ctx, "User not found")
        return nil, appstatus.NotFound("user not found.")
    }

    return user, nil
}

func (uc *EditAccountUsecase) checkPassword(ctx context.Context, user *entities.UserEntity, password string) error {
    if passwordutil.HashPassword(password) != user.PasswordHash.ValueOrZero() {
        return appstatus.InvalidCredential("Wrong password")
    }

    return nil
}

func (uc *EditAccountUsecase) updateEmail(user *entities.UserEntity, email string) error {
    if user == nil {
        return errors.New("user entity is nil")
    }

    if email == user.Email.ValueOrZero() {
        return nil
    }

    if err := validator.ValidateEmail(email); err != nil {
        return appstatus.BadEmail(err.Error())
    }

    user.Email.Set(email)

    return nil
}

func (uc *EditAccountUsecase) updatePassword(user *entities.UserEntity, newPassword string) error {
    if user == nil {
        return errors.New("user entity is nil")
    }

    newPasswordHash := passwordutil.HashPassword(newPassword)
    if newPasswordHash == user.PasswordHash.ValueOrZero() {
        return nil
    }

    if err := validator.ValidatePassword(newPassword); err != nil {
        return appstatus.BadPassword(err.Error())
    }

    user.PasswordHash.Set(newPasswordHash)

    return nil
}

func (uc *EditAccountUsecase) updateDiscoverable(user *entities.UserEntity, discoverable bool) error {
    if user == nil {
        return errors.New("user entity is nil")
    }

    user.Discoverable = discoverable

    return nil
}
