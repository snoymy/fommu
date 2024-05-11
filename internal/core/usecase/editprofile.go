package usecase

import (
	"app/internal/appstatus"
	"app/internal/core/dto"
	"app/internal/core/repo"
	"app/internal/core/validator"
	"context"
	"time"
)

type EditProfileUsecase struct {
    userRepo repo.UsersRepo
}

func NewEditProfileUsecase(userRepo repo.UsersRepo) *EditProfileUsecase {
    return &EditProfileUsecase{
        userRepo: userRepo,
    }
}

func (uc *EditProfileUsecase) Exec(ctx context.Context, username string, profile dto.UserProfileDTO) error {
    profileUpdate := false

    if username == "" {
        return appstatus.BadValue("username is empty.")
    }

    user, err := uc.userRepo.FindUserByUsername(ctx, username)
    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound("user not found.")
    }

    if !profile.Displayname.IsNull() {
        displayname := profile.Displayname.ValueOrZero()
        if displayname == "" {
            displayname = user.Username
        } else { 
            if err := validator.ValidateDisplayname(displayname); err != nil {
                return err
            }
        }
        user.Displayname = displayname
        profileUpdate = true
    }

    if !profile.NamePrefix.IsNull() {
        user.NamePrefix = profile.NamePrefix
        profileUpdate = true
    }

    if !profile.NameSuffix.IsNull() {
        user.NamePrefix = profile.NamePrefix
        profileUpdate = true
    }

    // if !profile.Gender.IsNull() {
    //     user.Gender = profile.Gender
    //     profileUpdate = true
    // }

    // if !profile.Pronoun.IsNull() {
    //     user.Pronoun = profile.Pronoun
    //     profileUpdate = true
    // }

    // if !profile.DateOfBirth.IsNull() {
    //     user.DateOfBirth = profile.DateOfBirth
    //     profileUpdate = true
    // }

    // if !profile.ShowDateOfBirth.IsNull() {
    //     user.ShowDateOfBirth = profile.ShowDateOfBirth.ValueOrZero()
    //     profileUpdate = true
    // }

    if !profile.Bio.IsNull() {
        user.Bio = profile.Bio
        profileUpdate = true
    }

    if !profile.Discoverable.IsNull() {
        user.Discoverable = profile.Discoverable.ValueOrZero()
        profileUpdate = true
    }

    if !profile.AutoApproveFollower.IsNull() {
        user.AutoApproveFollower = profile.AutoApproveFollower.ValueOrZero()
        profileUpdate = true
    }

    if !profile.Avatar.IsNull() {
        user.Avatar = profile.Avatar
        profileUpdate = true
    }

    if !profile.Banner.IsNull() {
        user.Banner = profile.Banner
        profileUpdate = true
    }

    if profileUpdate {
        user.UpdateAt.Set(time.Now().UTC())
        if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
            return err
        }
    }

    return nil
}
