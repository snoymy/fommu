package controller

import (
	"app/internal/appstatus"
	"app/internal/core/dto"
	"app/internal/core/usecase"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type UsersController struct {
    signup *usecase.SignupUsecase
    getUser *usecase.GetUserUsecase
    editProfile *usecase.EditProfileUsecase
}

func NewUsersController(signup *usecase.SignupUsecase, getUser *usecase.GetUserUsecase, editProfile *usecase.EditProfileUsecase) *UsersController {
    return &UsersController{
        signup: signup,
        getUser: getUser,
        editProfile: editProfile,
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

func (c *UsersController) EditProfile(w http.ResponseWriter, r *http.Request) error {
    username := chi.URLParam(r, "username")

    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return appstatus.BadValue("Cannot decode json.")
    }
    
    profile := dto.UserProfileDTO{}
    if body["displayname"] != nil {
        if v, ok := body["displayname"].(string); ok {
            profile.Displayname.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["namePrefix"] != nil {
        if v, ok := body["namePrefix"].(string); ok {
            profile.NamePrefix.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["nameSuffix"] != nil {
        if v, ok := body["nameSuffix"].(string); ok {
            profile.NameSuffix.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["gender"] != nil {
        if v, ok := body["gender"].(string); ok {
            profile.Gender.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["pronoun"] != nil {
        if v, ok := body["pronoun"].(string); ok {
            profile.Pronoun.Set(v)    
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
            profile.DateOfBirth.Set(t)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["showDateOfBirth"] != nil {
        if v, ok := body["showDateOfBirth"].(bool); ok {
            profile.ShowDateOfBirth.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["bio"] != nil {
        if v, ok := body["bio"].(string); ok {
            profile.Bio.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["discoverable"] != nil {
        if v, ok := body["discoverable"].(bool); ok {
            profile.Discoverable.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["autoApproveFollower"] != nil {
        if v, ok := body["autoApproveFollower"].(bool); ok {
            profile.ShowDateOfBirth.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["avatar"] != nil {
        if v, ok := body["avatar"].(string); ok {
            profile.Bio.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if body["banner"] != nil {
        if v, ok := body["banner"].(string); ok {
            profile.Bio.Set(v)    
        } else {
            return appstatus.BadValue("Invalid value")
        }
    }

    if err := c.editProfile.Exec(r.Context(), username, profile); err != nil {
        return err
    }

    return nil
}
