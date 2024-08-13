package controller

import (
	"app/internal/adapter/mapper"
	"app/internal/application/activitypub/usecase"
	"app/internal/core/appstatus"
	"app/internal/log"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/go-ap/jsonld"
	"github.com/go-chi/chi/v5"
	"github.com/snoymy/activitypub"
)

type APUsersController struct {
    getUser         *usecase.GetUserUsecase         `injectable:""`
    followUser      *usecase.ProcessFollowActivityUsecase      `injectable:""`
    createActivity  *usecase.CreateActivityUsecase  `injectable:""`
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

    person, err := mapper.UserToPerson(user)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    log.Info(ctx, "Build JSON-LD")
    bytes, err := jsonld.WithContext(
        jsonld.IRI(activitypub.ActivityBaseURI),
    ).Marshal(person)

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

    activity, err := mapper.JsonToActivity(string(body))
    if err := f.createActivity.Exec(r.Context(), activity); err != nil {
        log.Info(r.Context(), "Response with Error: " + err.Error())
        return err
    }

    // switch activity.Type {
    //     case activitypub.FollowType: 
    //         f.followUser.Exec(r.Context(), activity)
    // }

    w.WriteHeader(http.StatusAccepted)

    return nil
}
