package controller

import (
	"app/internal/activitypub/core/usecase"
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/log"
	"app/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-ap/jsonld"
	"github.com/go-chi/chi/v5"
	"github.com/snoymy/activitypub"
)

type APUsersController struct {
    getUser    *usecase.GetUserUsecase    `injectable:""`
    followUser *usecase.FollowUserUsecase `injectable:""`
}

func NewAPUsersController() *APUsersController {
    return &APUsersController{}
}

func (f *APUsersController) GetUser(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Get Accept type from header")
    acceptHeader := r.Header.Get("Accept")
    if !strings.Contains(acceptHeader, "application/activity+json") && !strings.Contains(acceptHeader, jsonld.ContentType) {
        err := appstatus.NotAccept(fmt.Sprintf("Invalid header: Accept is %s, not activity type.", acceptHeader))
        log.Info(ctx, "Response with Error: " + err.Error())
        return err
    }

    log.Info(ctx, "Get username from path param")
    username := chi.URLParam(r, "username")

    user, err := f.getUser.Exec(ctx, username)
    if err != nil {
        log.Error(ctx, "Response with Error: " + err.Error())
        return err
    }

    if user == nil {
        log.Info(ctx, "User not found")
        return appstatus.NotFound()
    }

    log.Info(ctx, "Build user url")
    userURL, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }
    inboxURL, err := url.JoinPath(userURL, "inbox")
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }
    outbox, err := url.JoinPath(userURL, "outbox")
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }
    followersURL, err := url.JoinPath(userURL, "followers")
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }
    followingURL, err := url.JoinPath(userURL, "following")
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    log.Info(ctx, "Create ActivityPub Person")
    p := activitypub.PersonNew(activitypub.IRI(userURL))

    p.Name = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, user.Displayname))
    p.PreferredUsername = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, user.Username))
    p.Inbox = activitypub.IRI(inboxURL)
    p.Outbox = activitypub.IRI(outbox)
    p.Followers = activitypub.IRI(followersURL)
    p.Following = activitypub.IRI(followingURL)
    p.PublicKey = activitypub.PublicKey{
        ID: activitypub.IRI(userURL + "#main-key"),
        Owner: activitypub.IRI(userURL),
        PublicKeyPem: user.PublicKey,
    }
    p.Summary = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(
        activitypub.DefaultLang, 
        strings.ReplaceAll(strings.ReplaceAll(utils.Linkify(user.Bio.ValueOrZero()), "\n", "<br>"), " ", "&nbsp;"),
    ))
    p.URL = activitypub.IRI(userURL)
    p.Icon = activitypub.Image{
        Type: activitypub.ImageType,
        MediaType: activitypub.MimeType(utils.GetMIMEFromExtension(filepath.Ext(user.Avatar.ValueOrZero()))),
        URL: activitypub.IRI(user.Avatar.ValueOrZero()),
    }
    p.Image = activitypub.Image{
        Type: activitypub.ImageType,
        MediaType: activitypub.MimeType(utils.GetMIMEFromExtension(filepath.Ext(user.Banner.ValueOrZero()))),
        URL: activitypub.IRI(user.Banner.ValueOrZero()),
    }
    p.Attachment = activitypub.ItemCollection{}
    for _, item := range user.Attachment.ValueOrZero() {
        attachment, err := utils.MapToStruct[activitypub.Object](item.(map[string]interface{}))
        if err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return appstatus.InternalServerError("Something went wrong.")
        }
        p.Attachment.Append(attachment)
    }
    p.Tag = activitypub.ItemCollection{}
    for _, item := range user.Tag.ValueOrZero() {
        tag, err := utils.MapToStruct[activitypub.Object](item.(map[string]interface{}))
        if err != nil {
            log.Error(ctx, "Error: " + err.Error())
            return appstatus.InternalServerError("Something went wrong.")
        }
        p.Tag.Append(tag)
    }

    log.Info(ctx, "Build JSON-LD")
    bytes, err := jsonld.WithContext(
        jsonld.IRI(activitypub.ActivityBaseURI),
    ).Marshal(p)

    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/activity+json")
    _, err = w.Write(bytes)

    return err
}

func (f *APUsersController) Inbox(w http.ResponseWriter, r *http.Request) error {
    acceptHeader := r.Header.Get("Signature")
    // j, _ := json.Marshal(r.Header)
    fmt.Println("Signature", acceptHeader)
    requestBuffer, err := httputil.DumpRequest(r, true)
    if err != nil {
    fmt.Println("Error dumping request:", err)
    return nil
  }

  // Print the raw request for debugging
  fmt.Println("Raw Request:")
  fmt.Println(string(requestBuffer))
    // if !strings.Contains(acceptHeader, "application/activity+json") && !strings.Contains(acceptHeader, jsonld.ContentType) {
    //     return appstatus.NotAccept(fmt.Sprintf("Invalid header: Accept is %s, not activity type.", acceptHeader))
    // }

    //username := chi.URLParam(r, "username")
    body, err := io.ReadAll(r.Body)
    if err != nil {
        return err
    }
    fmt.Printf("\n\nInbox: %s\n", body)

    activity := &activitypub.Activity{}
    err = json.Unmarshal(body, &activity)
    if err != nil {
        return err
    }

    var tagScheme struct {
        Tag []*activitypub.Object `json:"tag"`
    }
    err = json.Unmarshal(body, &tagScheme)
    if err != nil {
        return err
    }

    for _, item := range tagScheme.Tag {
        activity.Tag.Append(item)
    }

    switch activity.Type {
        case activitypub.FollowType: 
            f.followUser.Exec(r.Context(), activity)
    }

    return nil
}
