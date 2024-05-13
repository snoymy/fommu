package repo

import (
	"app/internal/api/core/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepoImpl(db *sqlx.DB) *UserRepository {
    return &UserRepository{
        db: db,
    }
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*entity.UserEntity, error) {
    var users []*entity.UserEntity = nil
    err := r.db.Select(&users, "select * from users where username=$1", username)
    if err != nil {
        return nil, err
    }

    if users == nil {
        return nil, nil
    }

    return users[0], nil
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
            bio, avatar, banner, tag, discoverable, auto_approve_follower, follower_count, following_count, 
            public_key, private_key, url, remote, redirect_url, create_at, update_at
        )
        values
        ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)
        `,
        user.ID, user.Email, user.PasswordHash, user.Status, user.Username, user.Displayname,
        user.NamePrefix, user.NameSuffix, user.Bio, user.Avatar, user.Banner, user.Tag, user.Discoverable, 
        user.AutoApproveFollower, user.FollowerCount, user.FollowingCount, user.PublicKey, user.PrivateKey, 
        user.URL, user.Remote, user.RedirectURL, user.CreateAt, user.UpdateAt,
    )

    if err != nil {
        return err
    }

    return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.UserEntity) error {
    _, err := r.db.Exec(
        `
        update users set 
            display_name=$1, name_prefix=$2, name_suffix=$3, bio=$4, 
            avatar=$5, banner=$6, discoverable=$7, auto_approve_follower=$8, update_at=$9
        where id = $10
        `,
        user.Displayname, user.NamePrefix, user.NameSuffix, user.Bio,
        user.Avatar, user.Banner, user.Discoverable, user.AutoApproveFollower, user.UpdateAt,
        user.ID,
    )

    if err != nil {
        return err
    }

    return nil
}
