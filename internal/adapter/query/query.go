package query

import (
	"app/internal/adapter/httpclient"
	"app/internal/core/entity"
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/snoymy/activitypub"
)

type Query struct {
    db *sqlx.DB `injectable:""`
    client httpclient.ActivitypubClient `injectable:""`
}

func NewQuery() *Query {
    return &Query{}
}

func (q *Query) FindUserById(ctx context.Context, id string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := q.db.Select(&users, "select * from users where id=$1", id)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (q *Query) SearchUser(ctx context.Context, textSearch string, domain string) ([]*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    textSearch = strings.ReplaceAll(textSearch, "%", "\\%")
    textSearch = strings.ReplaceAll(textSearch, "_", "\\_")
    err := q.db.Select(&users, "select * from users where (trim($1) <> '' and username ilike $1 || '%') and (trim($2) = '' or domain ilike $2 || '%') or (trim($1) <> '' and display_name ilike $1 || '%') limit 10", textSearch, domain)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users, nil
}

func (q *Query) FindPersonByActorId(ctx context.Context, url string) (*activitypub.Person, error) {
    actor, err := q.client.FetchActor(ctx, url)
    if err != nil {
        return nil, err
    }

    if actor == nil {
        return nil, nil
    }

    person := (*activitypub.Person)(actor)

    return person, nil
}

func (q *Query) FindUserByUsername(ctx context.Context, username string, domain string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := q.db.Select(&users, "select * from users where username=$1 and domain=$2", username, domain)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (q *Query) FindUserByActorId(ctx context.Context, actorId string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := q.db.Select(&users, "select * from users where actor_id=$1", actorId)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (q *Query) FindUserByEmail(ctx context.Context, email string, domain string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := q.db.Select(&users, "select * from users where email=$1 and domain=$2", email, domain)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (q *Query) FindUserByResourceName(ctx context.Context, resource string, domain string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := q.db.Select(&users, "select * from users where username||'@'||$1=$2", domain, resource)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
}

func (q *Query) FindPersonByUsername(ctx context.Context, username string, domain string) (*activitypub.Person, error) {
    links, err := q.client.FetchWebfinger(ctx, username, domain)
    if err != nil {
        return nil, err
    }

    href := ""
    for _, l := range links {
        var ok bool

        link, ok := l.(map[string]interface{})
        if !ok {
            continue
        }

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

        href, ok = link["href"].(string)
        if !ok {
            continue
        }

        break
    }

    if href == "" {
        return nil, nil
    }

    actor, err := q.client.FetchActor(ctx, href)
    if err != nil {
        return nil, err
    }

    if actor == nil {
        return nil, nil
    }

    person := (*activitypub.Person)(actor)

    return person, nil
}

func (q *Query) FindPersonFollowers(ctx context.Context, person *activitypub.Person, page int) (*activitypub.OrderedCollectionPage, error) {
    orderedCollection, err := q.client.FetchOrderedCollectionPage(ctx, person.Followers.GetLink().String(), page)
    if err != nil {
        return nil, err
    }

    return orderedCollection, nil
}

func (q *Query) FindPersonFollowing(ctx context.Context, person *activitypub.Person, page int) (*activitypub.OrderedCollectionPage, error) {
    orderedCollection, err := q.client.FetchOrderedCollectionPage(ctx, person.Following.GetLink().String(), page)
    if err != nil {
        return nil, err
    }

    return orderedCollection, nil
}

func (q *Query) FindSessionById(ctx context.Context, id string) (*entity.SessionEntity, error) {
    var sessions []*entity.SessionEntity = nil
    err := q.db.Select(&sessions, "select * from sessions where id=$1", id)

    if err != nil {
        return nil, err
    }

    if sessions == nil {
        return nil, nil
    }

    return sessions[0], nil
}

func (q *Query) FindActivityByActivityId(ctx context.Context, activityId string) (*entity.ActivityEntity, error) {
    var activity []*entity.ActivityEntity = nil
    err := q.db.Select(&activity, "select * from activities where activity->>'id'::text = $1", activityId)

    if err != nil {
        return nil, err
    }

    if activity == nil {
        return nil, nil
    }

    return activity[0], nil
}

func (q *Query) FindActivityById(ctx context.Context, activityId string) (*entity.ActivityEntity, error) {
    var activity []*entity.ActivityEntity = nil
    err := q.db.Select(&activity, "select * from activities where id = $1", activityId)

    if err != nil {
        return nil, err
    }

    if activity == nil {
        return nil, nil
    }

    return activity[0], nil
}

func (q *Query) CountFollows(ctx context.Context, follow *entity.FollowEntity) (int, error) {
    var rowsCount []int = nil
    err := q.db.Select(&rowsCount, "select count(*) from follows where follower = $1 and following = $2", follow.Follower, follow.Following)
    if err != nil {
        return 0, err
    }

    if rowsCount == nil {
        return 0, nil
    }

    return rowsCount[0], nil
}

