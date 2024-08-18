package repoimpls

import (
	"app/internal/adapter/commands"
	"app/internal/adapter/mappers"
	"app/internal/adapter/queries"
	"app/internal/config"
	"app/internal/core/entities"
	"context"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type UserRepoImpl struct {
    queries *queries.Query `injectable:""`
    commands *commands.Command `injectable:""`
}

func NewUserRepoImpl() *UserRepoImpl {
    return &UserRepoImpl{}
}

func (r *UserRepoImpl) FindUserByID(ctx context.Context, id string) (*entities.UserEntity, error) {
    user, err := r.queries.FindUserById(ctx, id)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (r *UserRepoImpl) FindUserByUsername(ctx context.Context, username string, domain string) (*entities.UserEntity, error) {
    if domain == "" {
        domain = config.Fommu.Domain
    }

    user, err := r.queries.FindUserByUsername(ctx, username, domain)
    if err != nil {
        return nil, err
    }

    if user != nil {
        if user.Domain == config.Fommu.Domain {
            return user, nil
        }

        if !user.UpdateAt.IsNull() && time.Now().UTC().Before(user.UpdateAt.ValueOrZero().Add(15 * time.Minute)) {
            return user, nil
        }

        person, err := r.queries.FindPersonByActorId(ctx, user.ActorId)
        if err != nil {
            return user, err
        }
        if person == nil {
            return user, nil
        }

        followers, err := r.queries.FindPersonFollowers(ctx, person, 0)
        if err != nil {
            return nil, err
        }

        following, err := r.queries.FindPersonFollowing(ctx, person, 0)
        if err != nil {
            return nil, err
        }
        
        userTemp, err := mappers.PersonToUser(person)

        if followers != nil {
            userTemp.FollowerCount = int(followers.TotalItems)
        }
        if following != nil {
            userTemp.FollowingCount = int(following.TotalItems)
        }

        hasUpdate := false
        if user.Displayname != userTemp.Displayname {
            user.Displayname = userTemp.Displayname 
            hasUpdate = true
        } 
        if user.Bio.ValueOrZero() != userTemp.Bio.ValueOrZero() {
            user.Bio = userTemp.Bio
            hasUpdate = true
        } 
        if user.FollowersURL != userTemp.FollowersURL {
            user.FollowersURL = userTemp.FollowersURL
            hasUpdate = true
        } 
        if user.FollowingURL != userTemp.FollowingURL {
            user.FollowingURL = userTemp.FollowingURL
            hasUpdate = true
        } 
        if user.InboxURL != userTemp.InboxURL {
            user.InboxURL = userTemp.InboxURL
            hasUpdate = true
        } 
        if user.OutboxURL != userTemp.OutboxURL {
            user.OutboxURL = userTemp.OutboxURL
            hasUpdate = true
        } 
        if user.Avatar != userTemp.Avatar {
            user.Avatar = userTemp.Avatar
            hasUpdate = true
        } 
        if user.Banner.ValueOrZero() != userTemp.Banner.ValueOrZero(){
            user.Banner = userTemp.Banner
            hasUpdate = true
        } 
        if reflect.DeepEqual(user.Attachment.ValueOrZero(), userTemp.Attachment.ValueOrZero()) {
            user.Attachment = userTemp.Attachment
            hasUpdate = true
        } 
        if reflect.DeepEqual(user.Tag.ValueOrZero(), userTemp.Tag.ValueOrZero()) {
            user.Tag = userTemp.Tag
            hasUpdate = true
        } 
        if user.FollowerCount != userTemp.FollowerCount {
            user.FollowerCount = userTemp.FollowerCount
            hasUpdate = true
        } 
        if user.FollowingCount != userTemp.FollowingCount {
            user.FollowingCount = userTemp.FollowingCount
            hasUpdate = true
        }

        if hasUpdate {
            user.UpdateAt.Set(time.Now().UTC())
            go r.UpdateUser(ctx, user)
        }

        return user, nil
    }

    person, err := r.queries.FindPersonByUsername(ctx, username, domain)

    if person == nil {
        return nil, nil
    }

    followers, err := r.queries.FindPersonFollowers(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    following, err := r.queries.FindPersonFollowing(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    user, err = mappers.PersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    user.Remote = true
    user.Discoverable = true
    user.CreateAt = time.Now().UTC()

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    go r.CreateUser(ctx, user)

    return user, nil
}

func (r *UserRepoImpl) FindUserByActorId(ctx context.Context, actorId string) (*entities.UserEntity, error) {
    user, err := r.queries.FindUserByActorId(ctx, actorId)
    if err != nil {
        return nil, err
    }
    
    if user != nil {
        return user, nil
    }

    person, err := r.queries.FindPersonByActorId(ctx, actorId)
    if err != nil {
        return nil, err
    }
    if person == nil {
        return nil, nil
    }

    followers, err := r.queries.FindPersonFollowers(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    following, err := r.queries.FindPersonFollowing(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    user, err = mappers.PersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    user.Remote = true
    user.Discoverable = true

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    go r.CreateUser(ctx, user)

    return user, nil
}

func (r *UserRepoImpl) FindResource(ctx context.Context, resource string, domain string) (*entities.UserEntity, error) {
    user, err := r.queries.FindUserByResourceName(ctx, resource, domain)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (r *UserRepoImpl) FindUserByEmail(ctx context.Context, email string, domain string) (*entities.UserEntity, error) {
    user, err := r.queries.FindUserByEmail(ctx, email, domain)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (r *UserRepoImpl) SearchUser(ctx context.Context, textSearch string, domain string) ([]*entities.UserEntity, error) {
    users, err := r.queries.SearchUser(ctx, textSearch, domain)
    if users != nil {
        return users, nil
    }

    person, err := r.queries.FindPersonByUsername(ctx, textSearch, domain)

    if person == nil {
        return nil, nil
    }

    followers, err := r.queries.FindPersonFollowers(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    following, err := r.queries.FindPersonFollowing(ctx, person, 0)
    if err != nil {
        return nil, err
    }

    user, err := mappers.PersonToUser(person)
    user.ID = uuid.New().String()
    user.Remote = true
    user.Remote = true
    user.Discoverable = true

    user.CreateAt = time.Now().UTC()

    if followers != nil {
        user.FollowerCount = int(followers.TotalItems)
    }
    if following != nil {
        user.FollowingCount = int(following.TotalItems)
    }

    go r.CreateUser(ctx, user)

    return []*entities.UserEntity{user}, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user *entities.UserEntity) error {
    err := r.commands.CreateUser(ctx, user)
    if err != nil {
        return err
    }

    return nil
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, user *entities.UserEntity) error {
    err := r.commands.UpdateUser(ctx, user)
    if err != nil {
        return err
    }

    return nil
}


