package controller

import (
	"app/internal/config"
	"app/internal/core/usecase"
	"encoding/json"
	"net/http"
	"net/url"
)

type WellKnown struct {
   findresource *usecase.FindResourceUsecase 
}

func NewWellKnownController(findactor *usecase.FindResourceUsecase) *WellKnown {
    return &WellKnown{
        findresource: findactor,
    }
}

func (f *WellKnown) WebFinger(w http.ResponseWriter, r *http.Request) error {
    resource := r.URL.Query().Get("resource")

    user, err := f.findresource.Exec(r.Context(), resource)
    if err != nil {
        return err
    }
    if user == nil {
        w.WriteHeader(404)
        return nil
    }

    userURL, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        return err
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

    w.Header().Add("Content-Type", "application/json")
    _, err = w.Write(bytes)

    return err
}
