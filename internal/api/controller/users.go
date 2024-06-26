package controller

import (
	"app/internal/api/core/dto"
	"app/internal/api/core/usecase"
	"app/internal/appstatus"
	"app/internal/types"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type UsersController struct {
    signup *usecase.SignupUsecase
    getUser *usecase.GetUserUsecase
    editProfile *usecase.EditProfileUsecase
    editAccount *usecase.EditAccountUsecase
    searchUser *usecase.SearchUserUsecase
}

func NewUsersController(
    signup *usecase.SignupUsecase, 
    getUser *usecase.GetUserUsecase, 
    editProfile *usecase.EditProfileUsecase,
    editAccount *usecase.EditAccountUsecase,
    searchUser *usecase.SearchUserUsecase,
) *UsersController {
    return &UsersController{
        signup: signup,
        getUser: getUser,
        editProfile: editProfile,
        editAccount: editAccount,
        searchUser: searchUser,
    }
}

func (c *UsersController) SignUp(w http.ResponseWriter, r *http.Request) error {
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return appstatus.BadValue("Cannot decode json.")
    }

    username, _ := body["username"].(string)
    password, _ := body["password"].(string)
    email, _ := body["email"].(string)

    err := c.signup.Exec(r.Context(), email, username, password)

    if err != nil {
        return err
    }

    return nil
}

func (c *UsersController) LookUp(w http.ResponseWriter, r *http.Request) error {
    username := r.URL.Query().Get("acct")
    
    user, err := c.getUser.Exec(r.Context(), username)

    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound()
    }

    return nil
}

func (c *UsersController) Search(w http.ResponseWriter, r *http.Request) error {
    username := r.URL.Query().Get("acct")
    
    users, err := c.searchUser.Exec(r.Context(), username)

    if err != nil {
        return err
    }

    if users == nil {
        return appstatus.NotFound()
    }

    res := []map[string]interface{}{}

    for _, user := range users {
        res = append(
            res,
            map[string]interface{}{
                "id": user.ID,
                "username": user.Username,
                "displayname": user.Displayname,
                "avatar": user.Avatar.ValueOrZero(),
                "banner": user.Banner.ValueOrZero(),
                "domain": user.Domain,
                "tag": user.Tag.ValueOrZero(),
            },
        )
    }

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *UsersController) GetUser(w http.ResponseWriter, r *http.Request) error {
    username := chi.URLParam(r, "username")

    username = strings.ReplaceAll(username, "%40", "@")
    username = strings.ReplaceAll(username, "%3A", ":")
    user, err := c.getUser.Exec(r.Context(), username)

    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound()
    }

    res := map[string]interface{}{
        "id": user.ID,
        "username": user.Username,
        "displayname": user.Displayname,
        "name_prefix": user.NamePrefix.ValueOrZero(),
        "name_suffix": user.NameSuffix.ValueOrZero(),
        "avatar": user.Avatar.ValueOrZero(),
        "banner": user.Banner.ValueOrZero(),
        "bio": user.Bio.ValueOrZero(),
        "domain": user.Domain,
        "preference": user.Preference.ValueOrZero(),
        "tag": user.Tag.ValueOrZero(),
        "follower_count": user.FollowerCount,
        "following_count": user.FollowingCount,
        "create_at": user.CreateAt.UTC(),
    }

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *UsersController) EditAccount(w http.ResponseWriter, r *http.Request) error {
    username := chi.URLParam(r, "username")

    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return appstatus.BadValue("Cannot decode json.")
    }
    
    account := dto.UserAccountDTO{}
    if i, ok := body["email"]; ok {
        if i == nil {
            value := types.Null[string]()
            account.Email = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            account.Email = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["currentPassword"]; ok {
        if i == nil {
            value := types.Null[string]()
            account.CurrentPassword = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            account.CurrentPassword = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["newPassword"]; ok {
        if i == nil {
            value := types.Null[string]()
            account.NewPassword = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            account.NewPassword = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if err := c.editAccount.Exec(r.Context(), username, account); err != nil {
        return err
    }

    return nil
}

func (c *UsersController) EditProfile(w http.ResponseWriter, r *http.Request) error {
    username := chi.URLParam(r, "username")

    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return appstatus.BadValue("Cannot decode json.")
    }
    
    profile := dto.UserProfileDTO{}
    if i, ok := body["displayname"]; ok {
        if i == nil {
            value := types.Null[string]()
            profile.Displayname = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            profile.Displayname = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["bio"]; ok {
        if i == nil {
            value := types.Null[string]()
            profile.Bio = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            profile.Bio = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["avatar"]; ok {
        if i == nil {
            value := types.Null[string]()
            profile.Avatar = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            profile.Avatar = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["banner"]; ok {
        if i == nil {
            value := types.Null[string]()
            profile.Banner = &value
        } else if v, ok := i.(string); ok {
            value := types.NewNullable(v)
            profile.Banner = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if i, ok := body["preference"]; ok {
        if i == nil {
            value := types.Null[types.JsonObject]()
            profile.Preference = &value
        } else if v, ok := i.(map[string]interface{}); ok {
            value := types.NewNullable(types.JsonObject(v))
            profile.Preference = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }
    /*
    if body["namePrefix"] != nil {
        if v, ok := body["namePrefix"].(string); ok {
            value := types.NewNullable[string](v)
            profile.NamePrefix = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["nameSuffix"] != nil {
        if v, ok := body["nameSuffix"].(string); ok {
            value := types.NewNullable[string](v)
            profile.NameSuffix = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["gender"] != nil {
        if v, ok := body["gender"].(string); ok {
            value := types.NewNullable[string](v)
            profile.Gender = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["pronoun"] != nil {
        if v, ok := body["pronoun"].(string); ok {
            value := types.NewNullable[string](v)
            profile.Pronoun = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["dateOfBirth"] != nil {
        if v, ok := body["dateOfBirth"].(string); ok {
            t, err := time.Parse(time.DateOnly, v)
            if err != nil {
                return err
            }
            value := types.NewNullable[time.Time](t)
            profile.DateOfBirth = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["showDateOfBirth"] != nil {
        if v, ok := body["showDateOfBirth"].(bool); ok {
            value := types.NewNullable[bool](v)
            profile.ShowDateOfBirth = &value   
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["discoverable"] != nil {
        if v, ok := body["discoverable"].(bool); ok {
            value := types.NewNullable[bool](v)
            profile.Discoverable = &value
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["autoApproveFollower"] != nil {
        if v, ok := body["autoApproveFollower"].(bool); ok {
            value := types.NewNullable[bool](v)
            profile.ShowDateOfBirth = &value   
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }
    */

    if err := c.editProfile.Exec(r.Context(), username, profile); err != nil {
        return err
    }

    return nil
}
