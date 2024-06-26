package usecase

import (
	"app/internal/api/core/dto"
	"app/internal/api/core/repo"
	"app/internal/api/core/validator"
	"app/internal/appstatus"
	"app/internal/config"
	"context"
	"html"
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

    if profile.Displayname != nil {
        displayname := profile.Displayname.ValueOrZero()
        if displayname == "" {
            displayname = user.Username
        } else { 
            if err := validator.ValidateDisplayname(displayname); err != nil {
                return err
            }
        }
        user.Displayname = html.EscapeString(displayname)
        hasUpdate = true
    }

    if profile.NamePrefix != nil {
        user.NamePrefix.Set(html.EscapeString(profile.NamePrefix.ValueOrZero()))
        hasUpdate = true
    }

    if profile.NameSuffix != nil {
        user.NameSuffix.Set(html.EscapeString(profile.NameSuffix.ValueOrZero()))
        hasUpdate = true
    }

    if profile.Preference != nil {
        user.Preference = *profile.Preference
        hasUpdate = true
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

    if profile.Bio != nil {
        user.Bio.Set(html.EscapeString(profile.Bio.ValueOrZero()))
        hasUpdate = true
    }

    if profile.Discoverable != nil {
        user.Discoverable = profile.Discoverable.ValueOrZero()
        hasUpdate = true
    }

    if profile.AutoApproveFollower != nil {
        user.AutoApproveFollower = profile.AutoApproveFollower.ValueOrZero()
        hasUpdate = true
    }

    if profile.Avatar != nil {
        user.Avatar = *profile.Avatar
        hasUpdate = true
    }

    if profile.Banner != nil {
        user.Banner = *profile.Banner
        hasUpdate = true
    }

    if hasUpdate {
        user.UpdateAt.Set(time.Now().UTC())
        if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
            return err
        }
    }

    return nil
}
