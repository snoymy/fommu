package repo

import (
	"app/internal/activitypub/core/entity"
	"app/internal/activitypub/core/repo"
	"app/internal/httpclient"
	"app/internal/types"
	"app/internal/utils"
	"context"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/snoymy/activitypub"
)

type UserRepository struct {
    db *sqlx.DB
    apClient *httpclient.ActivitypubClient
}

func NewUserRepoImpl(db *sqlx.DB, apClient *httpclient.ActivitypubClient) repo.UsersRepo {
    return &UserRepository{
        db: db,
        apClient: apClient,
    }
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string, domain string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := r.db.Select(&users, "select * from users where username=$1 and domain=$2", username, domain)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}
func (r *UserRepository) FindUserByActorId(ctx context.Context, actorId string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := r.db.Select(&users, "select * from users where actor_id=$1", actorId)
    if err != nil {
        return nil, err
    }

    if users != nil {
        return users[0], nil
    }

    println(actorId)

    person, err := r.apClient.GetUserByUrl(ctx, actorId)
    if err != nil {
        return nil, err
    }
    if person == nil {
        return nil, nil
    }

    followers, err := r.apClient.GetFollowersByUrl(ctx, person.Followers.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    following, err := r.apClient.GetFollowingByUrl(ctx, person.Following.GetLink().String(), 0)
    if err != nil {
        return nil, err
    }

    user, err := r.mapPersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    parsedUrl, err := url.Parse(user.ActorId.ValueOrZero())
    user.Domain = strings.TrimPrefix(parsedUrl.Hostname(), "www.")
    user.Remote = true
    user.Discoverable = true

    p := bluemonday.UGCPolicy()
    user.Username = p.Sanitize(user.Username)
    user.Displayname = p.Sanitize(user.Displayname)
    user.Bio.Set(p.Sanitize(user.Bio.ValueOrZero()))

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    go r.CreateUser(ctx, user)

    return user, nil
}

func (r *UserRepository) FindResource(ctx context.Context, resource string, domain string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := r.db.Select(&users, "select * from users where username||'@'||$1=$2", domain, resource)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := r.db.Select(&users, "select * from users where email=$1", email)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.UserEntity) error {
    _, err := r.db.Exec(
        `
        insert into users (
            id, email, password_hash, status, username, display_name, name_prefix, name_suffix, 
            bio, avatar, banner, attachment, tag, discoverable, auto_approve_follower, follower_count, following_count, 
            public_key, private_key, actor_id, url, inbox_url, outbox_url, followers_url, following_url, Domain, remote, redirect_url, 
            create_at, update_at
        )
        values
        ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
        `,
        user.ID, user.Email, user.PasswordHash, user.Status, user.Username, user.Displayname,
        user.NamePrefix, user.NameSuffix, user.Bio, user.Avatar, user.Banner, user.Attachment, user.Tag, user.Discoverable, 
        user.AutoApproveFollower, user.FollowerCount, user.FollowingCount, user.PublicKey, user.PrivateKey, user.ActorId,
        user.URL, user.InboxURL, user.OutboxURL, user.FollowersURL, user.FollowingURL, user.Domain, user.Remote, user.RedirectURL, 
        user.CreateAt, user.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (a *UserRepository) mapPersonToUser(person *activitypub.Person) (*entity.UserEntity, error) {
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

func (a *UserRepository) parseAttachment(person *activitypub.Person) (types.JsonArray, error) {
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

func (a *UserRepository) parseTag(person *activitypub.Person) (types.JsonArray, error) {
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
