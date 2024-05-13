package repo

import (
	"app/internal/activitypub/core/entity"
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
