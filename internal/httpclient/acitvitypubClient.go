package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-ap/activitypub"
)

type ActivitypubClient struct {

}

func NewActivitypubClient() *ActivitypubClient {
    return &ActivitypubClient{}
}

func (a *ActivitypubClient) FindUserByUsername(ctx context.Context, username string, domain string) (*activitypub.Person, error) {
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

    var person *activitypub.Person = nil
    for _, l := range links {
        link, ok := l.(map[string]interface{})
        rel, ok := link["rel"].(string)
        if !ok {
            continue
        }
        if rel != "self" {
            continue
        }

        hrefType, ok := link["type"].(string)
        if !ok {
            continue
        }
        if !strings.Contains(hrefType, "application/activity+json") && !strings.Contains(hrefType, "application/ld+json") {
            continue
        }

        href, ok := link["href"].(string)
        if !ok {
            continue
        }

        req, err := http.NewRequest("GET", href, nil)
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

        body, err = io.ReadAll(res.Body)
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

        break
    }

    if person == nil {
        return nil, nil
    }

    return person, nil
}

func (a *ActivitypubClient) GetUserByUrl(ctx context.Context, userUrl string) (*activitypub.Person, error) {
    var person *activitypub.Person = nil
    req, err := http.NewRequest("GET", userUrl, nil)
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

    if person == nil {
        return nil, nil
    }

    return person, nil
}

func (a *ActivitypubClient) GetFollowersByUrl(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
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
    followers := &activitypub.OrderedCollectionPage{}
    err = json.Unmarshal(body, &followers)
    if err != nil {
        return nil, err
    }

    return followers, nil
}

func (a *ActivitypubClient) GetFollowingByUrl(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
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
    
    following := &activitypub.OrderedCollectionPage{}
    err = json.Unmarshal(body, &following)
    if err != nil {
        return nil, err
    }

    return following, nil
}
