package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/snoymy/activitypub"
)

type ActivitypubClientImpl struct { }

func NewActivitypubClientImpl() ActivitypubClient {
    return &ActivitypubClientImpl{}
}

func (c *ActivitypubClientImpl) FetchWebfinger(ctx context.Context, domain string, username string) ([]interface{}, error) {
    urls := []string{
        fmt.Sprintf("https://%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("https://www.%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("http://%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("http://www.%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
    }

    var body []byte = nil
    for _, url := range urls {
        res, err := http.Get(url)
        if err != nil {
            continue
        }
        if res.StatusCode != http.StatusOK {
            continue
        }

        body, err = io.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }
        break
    }

    if body == nil {
        return nil, nil
    }

    var info map[string]interface{}
    err := json.Unmarshal(body, &info)
    if err != nil {
        return nil, err
    }

    links, ok := info["links"].([]interface{})
    if !ok {
        return nil, err
    }

    return links, nil
}

func (c *ActivitypubClientImpl) FetchActor(ctx context.Context, url string) (*activitypub.Actor, error) {
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

    actor := &activitypub.Actor{}

    err = json.Unmarshal(body, &actor)
    if err != nil {
        return nil, err
    }

    var nestedScheme struct {
        Tag []*activitypub.Object `json:"tag"`
        Attachment []*activitypub.Object `json:"attachment"`
    }

    err = json.Unmarshal(body, &nestedScheme)
    if err != nil {
        return nil, err
    }

    for _, item := range nestedScheme.Tag {
        actor.Tag.Append(item)
    }

    for _, item := range nestedScheme.Attachment {
        actor.Attachment.Append(item)
    }

    return actor, nil
}

func (c *ActivitypubClientImpl) FetchOrderedCollectionPage(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
    queryString := ""
    if page > 0 {
        queryString = fmt.Sprintf("?page=%d", page)
    }
    url = url + queryString

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
    
    //var temp types.JsonObject
    collectionPage := &activitypub.OrderedCollectionPage{}
    err = json.Unmarshal(body, &collectionPage)
    if err != nil {
        return nil, err
    }

    return collectionPage, nil
}
