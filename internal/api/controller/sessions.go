package controller

import (
	"app/internal/api/core/usecase"
	"app/internal/appstatus"
	"app/internal/types"
	"app/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type SessionsController struct {
    signin *usecase.SigninUsecase
    signout *usecase.SignOutUsecase
    refreshToken *usecase.RefreshTokenUsecase
    getToken *usecase.GetTokenUsecase
    revokeSession *usecase.RevokeSessionUsecase
    verifySession *usecase.GetIdentityUsecase
}

func NewSessionsController(
    signin *usecase.SigninUsecase, 
    signout *usecase.SignOutUsecase,
    refreshToken *usecase.RefreshTokenUsecase,
    getToken *usecase.GetTokenUsecase,
    revokeSession *usecase.RevokeSessionUsecase,
    verifySession *usecase.GetIdentityUsecase,
) *SessionsController {
    return &SessionsController{
        signin: signin,
        signout: signout,
        refreshToken: refreshToken,
        getToken: getToken,
        revokeSession: revokeSession,
        verifySession: verifySession,
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
        SameSite: http.SameSiteNoneMode,
        Secure: true,
        Path: "/",
    })

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) SignOut(w http.ResponseWriter, r *http.Request) error {
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        return appstatus.InvalidSession("Session not found")
    }

    if err := c.signout.Exec(r.Context(), sessionId.Value); err != nil {
        return err
    }

    http.SetCookie(w, &http.Cookie{
        Name: "session_id",
        Value: "", 
        MaxAge: -1,
        Expires: time.Unix(0, 0),
        HttpOnly: true,
        SameSite: http.SameSiteNoneMode,
        Secure: true,
        Path: "/",
    })

    return nil
}

func (c *SessionsController) RefreshToken(w http.ResponseWriter, r *http.Request) error {
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        return err
    }

    sessionId, err := r.Cookie("session_id")
    if err != nil {
        return appstatus.InvalidSession("Session not found")
    }

    refreshToken, _ := body["refresh_token"].(string)

    session, err := c.refreshToken.Exec(r.Context(), sessionId.Value, refreshToken)
    if err != nil {
        return err
    }

    res := map[string]interface{}{
        "access_token": session.AccessToken,
        "access_expire_at": session.AccessExpireAt.Unix(),
        "refresh_token": session.RefreshToken,
        "refresh_expire_at": session.RefreshExpireAt.Unix(),
    }

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) GetToken(w http.ResponseWriter, r *http.Request) error {
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        return appstatus.InvalidSession("Session not found")
    }

    session, err := c.getToken.Exec(r.Context(), sessionId.Value)
    if err != nil {
        return err
    }

    res := map[string]interface{}{
        "access_token": session.AccessToken,
        "access_expire_at": session.AccessExpireAt.Unix(),
        "refresh_token": session.RefreshToken,
        "refresh_expire_at": session.RefreshExpireAt.Unix(),
    }

    bytes, err := json.Marshal(res)

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) VerifySession(w http.ResponseWriter, r *http.Request) error {
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        return appstatus.InvalidSession("Session not found")
    }

    user, err := c.verifySession.Exec(r.Context(), sessionId.Value)

    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound()
    }

    res := map[string]interface{}{
        "id": user.ID,
        "email": user.Email,
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

func (c *SessionsController) RevokeSession(w http.ResponseWriter, r *http.Request) error {
    currentSessionId, err := r.Cookie("session_id")
    if err != nil {
        return appstatus.InvalidSession("Session not found")
    }

    sessionId := chi.URLParam(r, "sessionId")

    if err := c.revokeSession.Exec(r.Context(), currentSessionId.Value, sessionId); err != nil {
        return err
    }

    return nil
}
