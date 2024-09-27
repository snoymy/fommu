package controllers

import (
	"app/internal/application/fommu/usecases"
	"app/internal/application/appstatus"
	"app/internal/core/types"
	"app/internal/log"
	"app/internal/utils/requestutil"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type SessionsController struct {
    signin        *usecases.SigninUsecase        `injectable:""`
    signout       *usecases.SignOutUsecase       `injectable:""`
    refreshToken  *usecases.RefreshTokenUsecase  `injectable:""`
    getToken      *usecases.GetTokenUsecase      `injectable:""`
    revokeSession *usecases.RevokeSessionUsecase `injectable:""`
    verifySession *usecases.GetIdentityUsecase   `injectable:""`
}

func NewSessionsController() *SessionsController {
    return &SessionsController{}
}

func (c *SessionsController) Signin(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.BadValue("Cannot decode json.")
    }

    email, _ := body["email"].(string)
    password, _ := body["password"].(string)

    device, os := requestutil.GetClientPlatform(r)
    clientData := types.JsonObject{
        "os": os,
        "device": device,
    }

    session, err := c.signin.Exec(ctx, email, password, clientData)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
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
    if err != nil {
        log.Error(ctx, err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) SignOut(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    sessionId, err := r.Cookie("session_id")
    if err != nil {
        log.Warn(ctx, "Response with error: " + err.Error())
        return appstatus.InvalidSession("Session not found")
    }

    if err := c.signout.Exec(ctx, sessionId.Value); err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
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
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Decode json")
    var body map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        log.Error(ctx, "Response with error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong")
    }

    log.Info(ctx, "Get session id from cookie")
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return appstatus.InvalidSession("Session not found")
    }

    refreshToken, _ := body["refresh_token"].(string)

    session, err := c.refreshToken.Exec(ctx, sessionId.Value, refreshToken)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    res := map[string]interface{}{
        "access_token": session.AccessToken,
        "access_expire_at": session.AccessExpireAt.Unix(),
        "refresh_token": session.RefreshToken,
        "refresh_expire_at": session.RefreshExpireAt.Unix(),
    }

    bytes, err := json.Marshal(res)
    if err != nil {
        log.Error(ctx, err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) GetToken(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Get session id from cookie")
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        log.Warn(ctx, "Response with error: " + err.Error())
        return appstatus.InvalidSession("Session not found")
    }

    session, err := c.getToken.Exec(ctx, sessionId.Value)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    res := map[string]interface{}{
        "access_token": session.AccessToken,
        "access_expire_at": session.AccessExpireAt.Unix(),
        "refresh_token": session.RefreshToken,
        "refresh_expire_at": session.RefreshExpireAt.Unix(),
    }

    bytes, err := json.Marshal(res)
    if err != nil {
        log.Error(ctx, err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) VerifySession(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Get session id from cookie")
    sessionId, err := r.Cookie("session_id")
    if err != nil {
        log.Warn(ctx, "Response with error: " + err.Error())
        return appstatus.InvalidSession("Session not found")
    }

    user, err := c.verifySession.Exec(ctx, sessionId.Value)
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    if user == nil {
        err := appstatus.NotFound()
        log.Info(ctx, "Response with error: " + err.Error())
        return err
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
        "attachment": user.Attachment.ValueOrZero(),
        "follower_count": user.FollowerCount,
        "following_count": user.FollowingCount,
        "create_at": user.CreateAt.UTC(),
    }

    bytes, err := json.Marshal(res)
    if err != nil {
        log.Error(ctx, err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    w.Write(bytes)

    return nil
}

func (c *SessionsController) RevokeSession(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Get session id from cookie")
    currentSessionId, err := r.Cookie("session_id")
    if err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return appstatus.InvalidSession("Session not found")
    }

    sessionId := chi.URLParam(r, "sessionId")

    if err := c.revokeSession.Exec(ctx, currentSessionId.Value, sessionId); err != nil {
        log.Info(ctx, "Response with error: " + err.Error())
        return err
    }

    return nil
}
