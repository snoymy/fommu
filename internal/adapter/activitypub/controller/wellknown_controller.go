package controller

import (
	"app/internal/application/activitypub/usecase"
	"app/internal/core/appstatus"
	"app/internal/config"
	"app/internal/log"
	"encoding/json"
	"net/http"
	"net/url"
)

type WellKnown struct {
   findresource *usecase.FindResourceUsecase `injectable:""`
}

func NewWellKnownController() *WellKnown {
    return &WellKnown{}
}

func (f *WellKnown) WebFinger(w http.ResponseWriter, r *http.Request) error {
    ctx := r.Context()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Info(ctx, "Get resource from query param")
    resource := r.URL.Query().Get("resource")

    user, err := f.findresource.Exec(ctx, resource)
    if err != nil {
        log.Info(ctx, "Response with Error: ", err.Error())
        return err
    }
    if user == nil {
        log.Info(ctx, "User not found")
        w.WriteHeader(404)
        return nil
    }

    log.Info(ctx, "Build response")
    userURL, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        log.Error(ctx, "Error: ", err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    res := map[string]interface{}{
        "subject": "acct:" + user.Username + "@" + config.Fommu.Domain,
        "links": []interface{}{
            map[string]interface{}{
                "type":  "application/activity+json",
                "rel": "self",
                "href": userURL,
            },
            map[string]interface{}{
                "type":  "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
                "rel": "self",
                "href": userURL,
            },
        },
    }

    bytes, err := json.Marshal(res)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Something went wrong.")
    }

    w.Header().Add("Content-Type", "application/json")
    _, err = w.Write(bytes)

    return err
}
