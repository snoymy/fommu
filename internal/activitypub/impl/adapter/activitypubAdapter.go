package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-ap/activitypub"
)

type ActivitypubAdapterImpl struct {

}

func NewActivitypubAdapterImpl() *ActivitypubAdapterImpl {
    return &ActivitypubAdapterImpl{}
}

func (a *ActivitypubAdapterImpl) GetUserByUrl(ctx context.Context, url string) (*activitypub.Person, error) {
    var person *activitypub.Person = nil
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Accept", "application/activity+json")
    client := &http.Client{
        Timeout: 30 * time.Second,
    }
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        return nil, nil 
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    person = &activitypub.Person{}
    err = json.Unmarshal(body, &person)
    if err != nil {
        return nil, err
    }

    var tagScheme struct {
        Tag []*activitypub.Object `json:"tag"`
    }
    err = json.Unmarshal(body, &tagScheme)
    if err != nil {
        return nil, err
    }

    for _, item := range tagScheme.Tag {
        person.Tag.Append(item)
    }

    return person, nil
}
