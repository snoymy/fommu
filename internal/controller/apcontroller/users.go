package apcontroller

import (
	"app/internal/appstatus"
	"app/internal/config"
	"app/internal/core/usecase"
	"net/http"
	"net/url"

	ap "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/go-chi/chi/v5"
)

type APUsersController struct {
    getUser *usecase.GetUserUsecase
}

func NewAPUsersController(getUser *usecase.GetUserUsecase) *APUsersController {
    return &APUsersController{
        getUser: getUser,
    }
}

func (f *APUsersController) GetUser(w http.ResponseWriter, r *http.Request) error {
    username := chi.URLParam(r, "username")

    println(username)

    user, err := f.getUser.Exec(r.Context(), username)
    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound()
    }

    userURL, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        return err
    }
    inboxURL, err := url.JoinPath(userURL, "inbox")
    if err != nil {
        return err
    }
    outbox, err := url.JoinPath(userURL, "outbox")
    if err != nil {
        return err
    }
    followersURL, err := url.JoinPath(userURL, "followers")
    if err != nil {
        return err
    }
    followingURL, err := url.JoinPath(userURL, "following")
    if err != nil {
        return err
    }
    p := ap.PersonNew(ap.IRI(userURL))

    p.Name              = ap.NaturalLanguageValuesNew(ap.LangRefValueNew(ap.DefaultLang, user.Displayname))
    p.PreferredUsername = ap.NaturalLanguageValuesNew(ap.LangRefValueNew(ap.DefaultLang, user.Username))
    p.Inbox             = ap.IRI(inboxURL)
    p.Outbox            = ap.IRI(outbox)
    p.Followers         = ap.IRI(followersURL)
    p.Following         = ap.IRI(followingURL)
    p.PublicKey         = ap.PublicKey{
                              ID: ap.IRI(userURL + "#main-key"),
                              Owner: ap.IRI(userURL),
                              PublicKeyPem: user.PublicKey,
                          }
    p.Summary           = ap.NaturalLanguageValuesNew(ap.LangRefValueNew(ap.DefaultLang, user.Bio.ValueOrZero()))
    p.URL               = ap.IRI(userURL)
    p.Icon              = ap.Image{
                              Type: ap.ImageType,
                              URL: ap.IRI(user.Avatar.ValueOrZero()),
                          }

    bytes, err := jsonld.WithContext(
        jsonld.IRI(ap.ActivityBaseURI),
    ).Marshal(p)

    if err != nil {
        return err
    }

    w.Header().Add("Content-Type", "application/activity+json")
    _, err = w.Write(bytes)

    return err
}
