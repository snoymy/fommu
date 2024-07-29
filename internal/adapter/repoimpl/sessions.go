package repoimpl

import (
	"app/internal/core/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type SessionsRepoImpl struct {
    db *sqlx.DB `injectable:""`
}

func NewSessionRepoImpl() *SessionsRepoImpl {
    return &SessionsRepoImpl{}
}

func (r *SessionsRepoImpl) CreateSession(ctx context.Context, session *entity.SessionEntity) error {
    _, err := r.db.Exec(
        `
        insert into sessions (
            id, access_token, access_expire_at, refresh_token, refresh_expire_at, 
            owner, metadata, login_at, last_refresh
        )
        values
        ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        `,
        session.ID, session.AccessToken, session.AccessExpireAt, session.RefreshToken, session.RefreshExpireAt,
        session.Owner, session.Metadata, session.LoginAt, session.LastRefresh,
    )

    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) UpdateSession(ctx context.Context, session *entity.SessionEntity) error {
    _, err := r.db.Exec(
        `
        update sessions 
        set access_token=$1, access_expire_at=$2, refresh_token=$3, refresh_expire_at=$4, last_refresh=$5
        where id=$6
        `,
        session.AccessToken, session.AccessExpireAt, session.RefreshToken, session.RefreshExpireAt, session.LastRefresh, session.ID,
    )

    if err != nil {
        return err
    }

    return nil
}

func (r *SessionsRepoImpl) FindSessionByID(ctx context.Context, id string) (*entity.SessionEntity, error) {
    var sessions []*entity.SessionEntity = nil
    err := r.db.Select(&sessions, "select * from sessions where id=$1", id)

    if err != nil {
        return nil, err
    }

    if sessions == nil {
        return nil, nil
    }

    return sessions[0], nil
}


func (r *SessionsRepoImpl) DeleteSession(ctx context.Context, id string) error {
    _, err := r.db.Exec("delete from sessions where id=$1", id)
    if err != nil {
        return err
    }

    return nil
}
