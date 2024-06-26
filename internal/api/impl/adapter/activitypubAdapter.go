package adapter

import (
	"app/internal/api/core/entity"
	"app/internal/types"
	"app/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/snoymy/activitypub"
	"github.com/google/uuid"
)

type ActivitypubAdapterImpl struct {

}

func NewActivitypubAdapterImpl() *ActivitypubAdapterImpl {
    return &ActivitypubAdapterImpl{}
}

func (a *ActivitypubAdapterImpl) FindUserByUsername(ctx context.Context, username string, domain string) (*entity.UserEntity, error) {
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

    followers, err := a.GetFollowersByUrl(ctx, person.Followers.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    following, err := a.GetFollowingByUrl(ctx, person.Following.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    user, err := a.mapPersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    parsedUrl, err := url.Parse(user.ActorId.ValueOrZero())
    user.Domain = strings.TrimPrefix(parsedUrl.Hostname(), "www.")
    user.Remote = true
    user.Discoverable = true
    user.CreateAt = time.Now().UTC()

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    return user, nil
}

func (a *ActivitypubAdapterImpl) GetUserByUrl(ctx context.Context, userUrl string) (*entity.UserEntity, error) {
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

    followers, err := a.GetFollowersByUrl(ctx, person.Followers.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    following, err := a.GetFollowingByUrl(ctx, person.Following.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    user, err := a.mapPersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    parsedUrl, err := url.Parse(user.ActorId.ValueOrZero())
    user.Domain = strings.TrimPrefix(parsedUrl.Hostname(), "www.")
    user.Remote = true
    user.Discoverable = true

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    return user, nil
}

func (a *ActivitypubAdapterImpl) GetFollowersByUrl(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
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

func (a *ActivitypubAdapterImpl) GetFollowingByUrl(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
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

func (a *ActivitypubAdapterImpl) mapPersonToUser(person *activitypub.Person) (*entity.UserEntity, error) {
    user := entity.NewUserEntity()

    if person.ID.IsValid() { 
        if person.ID.IsLink() { 
            user.ActorId.Set(person.ID.GetLink().String())
        }
    }
    if person.URL != nil { 
        user.URL.Set(person.URL.GetLink().String())
    } else {
        if person.ID.IsValid() { 
            if person.ID.IsLink() { 
                user.URL.Set(person.ID.GetLink().String())
            }
        }
    }
    if person.PreferredUsername != nil { 
        user.Username = person.PreferredUsername.String() 
    }
    if person.Name != nil {
        user.Displayname = person.Name.String()
    }
    if person.Summary != nil {
        user.Bio.Set(person.Summary.String())
    }
    if person.Followers != nil {
        user.FollowersURL.Set(person.Followers.GetLink().String())
    }
    if person.Following != nil {
        user.FollowingURL.Set(person.Following.GetLink().String())
    }
    if person.Inbox != nil {
        user.InboxURL.Set(person.Inbox.GetLink().String())
    }
    if person.Outbox != nil {
        user.OutboxURL.Set(person.Outbox.GetLink().String())
    }
    if person.Icon != nil {
        user.Avatar.Set(person.Icon.(*activitypub.Image).URL.GetLink().String())
    }
    if person.Image != nil {
        user.Banner.Set(person.Image.(*activitypub.Image).URL.GetLink().String())
    }
    attachments, err := a.parseAttachment(person)
    if err != nil {
        return nil, err
    }
    user.Attachment.Set(attachments)
    tags, err := a.parseTag(person)
    if err != nil {
        return nil, err
    }
    user.Tag.Set(tags)
    user.PublicKey = person.PublicKey.PublicKeyPem

    return user, nil
}

func (a *ActivitypubAdapterImpl) parseAttachment(person *activitypub.Person) (types.JsonArray, error) {
    var err error
    attachments := types.JsonArray{}
    if person.Tag != nil {
        for _, item := range person.Attachment {
            var attachment interface{}
            if item.IsObject() {
                attachment, err = utils.StructToMap(item.(*activitypub.Object))
                if err != nil {
                    return nil, err
                }
            } else if item.IsLink() {
                attachment = item.(*activitypub.Link).GetLink().String()
            }
            attachments = append(attachments, attachment)
        }
    }

    return attachments, nil
}

func (a *ActivitypubAdapterImpl) parseTag(person *activitypub.Person) (types.JsonArray, error) {
    var err error
    tags := types.JsonArray{}
    if person.Tag != nil {
        for _, item := range person.Tag {
            var tag interface{}
            if item.IsObject() {
                tag, err = utils.StructToMap(item.(*activitypub.Object))
                if err != nil {
                    return nil, err
                }
            } else if item.IsLink() {
                tag = item.(*activitypub.Link).GetLink().String()
            }
            tags = append(tags, tag)
        }
    }

    return tags, nil
}
