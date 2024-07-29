package mapper

import (
	"app/internal/config"
	"app/internal/core/entity"
	"app/internal/core/types"
	"app/internal/utils/mimeutil"
	"app/internal/utils/stringutil"
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/snoymy/activitypub"
)

func PersonToUser(person *activitypub.Person) (*entity.UserEntity, error) {
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
    attachments, err := parseAttachment(person)
    if err != nil {
        return nil, err
    }
    user.Attachment.Set(attachments)
    tags, err := parseTag(person)
    if err != nil {
        return nil, err
    }
    user.Tag.Set(tags)
    user.PublicKey = person.PublicKey.PublicKeyPem

    return user, nil
}

func UserToPerson(user *entity.UserEntity) (*activitypub.Person, error) {
    if user == nil {
        return nil, errors.New("user is nil")
    }

    userURL, err := url.JoinPath(config.Fommu.URL, "users", user.Username)
    if err != nil {
        return nil, err
    }
    inboxURL, err := url.JoinPath(userURL, "inbox")
    if err != nil {
        return nil, err
    }
    outbox, err := url.JoinPath(userURL, "outbox")
    if err != nil {
        return nil, err
    }
    followersURL, err := url.JoinPath(userURL, "followers")
    if err != nil {
        return nil, err
    }
    followingURL, err := url.JoinPath(userURL, "following")
    if err != nil {
        return nil, err
    }

    person := activitypub.PersonNew(activitypub.IRI(userURL))

    person.Name = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, user.Displayname))
    person.PreferredUsername = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, user.Username))
    person.Inbox = activitypub.IRI(inboxURL)
    person.Outbox = activitypub.IRI(outbox)
    person.Followers = activitypub.IRI(followersURL)
    person.Following = activitypub.IRI(followingURL)
    person.PublicKey = activitypub.PublicKey{
        ID: activitypub.IRI(userURL + "#main-key"),
        Owner: activitypub.IRI(userURL),
        PublicKeyPem: user.PublicKey,
    }
    person.Summary = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(
        activitypub.DefaultLang, 
        strings.ReplaceAll(strings.ReplaceAll(stringutil.Linkify(user.Bio.ValueOrZero()), "\n", "<br>"), " ", "&nbsp;"),
    ))
    person.URL = activitypub.IRI(userURL)
    person.Icon = activitypub.Image{
        Type: activitypub.ImageType,
        MediaType: activitypub.MimeType(mimeutil.GetMIMEFromExtension(filepath.Ext(user.Avatar.ValueOrZero()))),
        URL: activitypub.IRI(user.Avatar.ValueOrZero()),
    }
    person.Image = activitypub.Image{
        Type: activitypub.ImageType,
        MediaType: activitypub.MimeType(mimeutil.GetMIMEFromExtension(filepath.Ext(user.Banner.ValueOrZero()))),
        URL: activitypub.IRI(user.Banner.ValueOrZero()),
    }
    person.Attachment = activitypub.ItemCollection{}
    for _, item := range user.Attachment.ValueOrZero() {
        attachment, err := MapToStruct[activitypub.Object](item.(map[string]interface{}))
        if err != nil {
            return nil, err
        }
        person.Attachment.Append(attachment)
    }
    person.Tag = activitypub.ItemCollection{}
    for _, item := range user.Tag.ValueOrZero() {
        tag, err := MapToStruct[activitypub.Object](item.(map[string]interface{}))
        if err != nil {
            return nil, err
        }
        person.Tag.Append(tag)
    }

    return person, nil
}

func parseAttachment(person *activitypub.Person) (types.JsonArray, error) {
    var err error
    attachments := types.JsonArray{}
    if person.Tag != nil {
        for _, item := range person.Attachment {
            var attachment interface{}
            if item.IsObject() {
                attachment, err = StructToMap(item.(*activitypub.Object))
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

func parseTag(person *activitypub.Person) (types.JsonArray, error) {
    var err error
    tags := types.JsonArray{}
    if person.Tag != nil {
        for _, item := range person.Tag {
            var tag interface{}
            if item.IsObject() {
                tag, err = StructToMap(item.(*activitypub.Object))
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
