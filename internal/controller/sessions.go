package controller

import (
	"app/internal/core/usecase"
	"app/internal/types"
	"app/internal/utils"
	"encoding/json"
	"net/http"
)

type SessionsController struct {
    signin *usecase.SigninUsecase
}

func NewSessionsController(signin *usecase.SigninUsecase) *SessionsController {
    return &SessionsController{
        signin: signin,
    }
}

func (c *SessionsController) Signin(w http.ResponseWriter, r *http.Request) error {
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return err
    }

    email, _ := body["email"].(string)
    password, _ := body["password"].(string)

    device, os := utils.GetClientPlatform(r)
    clientData := types.JsonObject{
        "os": os,
        "device": device,
    }

    session, err := c.signin.Exec(r.Context(), email, password, clientData)
    if err != nil {
        return err
    }

    res := map[string]interface{}{
        "access_token": session.AccessToken,
        "access_expire_at": session.AccessExpireAt.Unix(),
        "refresh_token": session.RefreshToken,
        "refresh_expire_at": session.RefreshExpireAt.Unix(),
    }

    http.SetCookie(w, &http.Cookie{
        Name: "session_id",
        Value: session.ID,
        MaxAge: int(session.RefreshExpireAt.Unix() - session.LoginAt.Unix()),
        Expires: session.RefreshExpireAt,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
        Secure: false,
        Path: "/",
    })

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}
