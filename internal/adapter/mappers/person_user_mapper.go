package mappers

import (
	"app/internal/adapter/model/activitypub_extended"
	"app/internal/core/entities"
	"app/internal/core/types"
	"app/internal/utils/mimeutil"
	"app/internal/utils/stringutil"
	"app/internal/utils/structutil"
	"errors"
	"html"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/snoymy/activitypub"
)

func PersonToUser(person *activitypub.Person) (*entities.UserEntity, error) {
    p := bluemonday.UGCPolicy()
    user := entities.NewUserEntity()

    if person.ID.IsValid() { 
        if person.ID.IsLink() { 
            user.ActorId = person.ID.GetLink().String()
        }
    }
    if person.URL != nil { 
        user.URL = person.URL.GetLink().String()
    } else {
        if person.ID.IsValid() { 
            if person.ID.IsLink() { 
                user.URL = person.ID.GetLink().String()
            }
        }
    }
    if person.PreferredUsername != nil { 
        user.Username = p.Sanitize(person.PreferredUsername.String())
    }
    if person.Name != nil {
        user.Displayname = p.Sanitize(person.Name.String())
    }
    if person.Summary != nil {
        user.Bio.Set(p.Sanitize(person.Summary.String()))
    }
    if person.Followers != nil {
        user.FollowersURL = person.Followers.GetLink().String()
    }
    if person.Following != nil {
        user.FollowingURL = person.Following.GetLink().String()
    }
    if person.Inbox != nil {
        user.InboxURL = person.Inbox.GetLink().String()
    }
    if person.Outbox != nil {
        user.OutboxURL = person.Outbox.GetLink().String()
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
    parsedUrl, err := url.Parse(user.ActorId)
    user.Domain = strings.TrimPrefix(parsedUrl.Hostname(), "www.")
    user.Tag.Set(tags)
    user.PublicKey = person.PublicKey.PublicKeyPem

    return user, nil
}

func UserToPerson(user *entities.UserEntity) (*activitypub.Person, error) {
    if user == nil {
        return nil, errors.New("user is nil")
    }

    person := activitypub.PersonNew(activitypub.IRI(user.ActorId))

    person.Name = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, html.UnescapeString(user.Displayname)))
    person.PreferredUsername = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(activitypub.DefaultLang, user.Username))
    person.Inbox = activitypub.IRI(user.InboxURL)
    person.Outbox = activitypub.IRI(user.OutboxURL)
    person.Followers = activitypub.IRI(user.FollowersURL)
    person.Following = activitypub.IRI(user.FollowingURL)
    person.PublicKey = activitypub.PublicKey{
        ID: activitypub.IRI(user.ActorId + "#main-key"),
        Owner: activitypub.IRI(user.ActorId),
        PublicKeyPem: user.PublicKey,
    }
    person.Summary = activitypub.NaturalLanguageValuesNew(activitypub.LangRefValueNew(
        activitypub.DefaultLang, 
        stringutil.Linkify(strings.ReplaceAll(strings.ReplaceAll(user.Bio.ValueOrZero(), "\n", "<br>"), " ", "&nbsp;")),
    ))
    person.URL = activitypub.IRI(user.URL)
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
        attachment, err := structutil.MapToStruct[activitypub_extended.PropertyValue](item.(map[string]interface{}))
        attachment.Name = activitypub.NaturalLanguageValuesNew(
            activitypub.LangRefValueNew(
                activitypub.DefaultLang, 
                stringutil.Linkify(strings.ReplaceAll(strings.ReplaceAll(attachment.Name.String(), "\n", "<br>"), " ", "&nbsp;")),
            ),
        )
        attachment.Value = activitypub.NaturalLanguageValuesNew(
            activitypub.LangRefValueNew(
                activitypub.DefaultLang, 
                stringutil.Linkify(strings.ReplaceAll(strings.ReplaceAll(attachment.Value.String(), "\n", "<br>"), " ", "&nbsp;")),
            ),
        )
        if err != nil {
            return nil, err
        }
        person.Attachment.Append(attachment)
    }
    person.Tag = activitypub.ItemCollection{}
    for _, item := range user.Tag.ValueOrZero() {
        tag, err := structutil.MapToStruct[activitypub.Object](item.(map[string]interface{}))
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
    if person.Attachment != nil {
        for _, item := range person.Attachment {
            var attachment interface{}
            if item.GetType() == activitypub_extended.PropertyValueType {
                attachment, err = structutil.StructToMap(item.(*activitypub_extended.PropertyValue))
                if err != nil {
                    return nil, err
                }
            } else if item.IsObject() {
                attachment, err = structutil.StructToMap(item.(*activitypub.Object))
                if err != nil {
                    return nil, err
                }
            } else if item.IsLink() {
                attachment = item.(*activitypub.Link).GetLink().String()
            } else {
                return nil, nil
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
                tag, err = structutil.StructToMap(item.(*activitypub.Object))
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
