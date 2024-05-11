CREATE TABLE IF NOT EXISTS users (
    id                      UUID PRIMARY KEY,
    email                   VARCHAR(255) NOT NULL UNIQUE,
    password_hash           VARCHAR(255) NOT NULL,
    status                  VARCHAR(255) NOT NULL, 
    username                VARCHAR(255) NOT NULL UNIQUE,
    display_name            VARCHAR(255),
    name_prefix             VARCHAR(255),
    name_suffix             VARCHAR(255),
    bio                     TEXT,
    avatar                  VARCHAR(255),
    banner                  VARCHAR(255),
    tag                     TEXT[],
    discoverable            BOOLEAN DEFAULT true,
    auto_approve_follower   BOOLEAN DEFAULT false,
    follower_count          INT DEFAULT 0,
    following_count         INT DEFAULT 0,
    public_key              TEXT NOT NULL,
    private_key             TEXT NOT NULL,
    url                     VARCHAR(255) NOT NULL,
    remote                  BOOLEAN DEFAULT false,
    redirect_url            VARCHAR(255),
    create_at               TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_at               TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sessions (
    id                  UUID PRIMARY KEY,
    access_token        TEXT NOT NULL,
    access_expire_at    TIMESTAMPTZ NOT NULL,
    refresh_token       TEXT NOT NULL,
    refresh_expire_at   TIMESTAMPTZ NOT NULL,
    owner               UUID REFERENCES users(id) NOT NULL,
    metadata            JSONB,
    login_at            TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_refresh        TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS media (
    id                  UUID PRIMARY KEY,
    url                 TEXT NOT NULL,
    type                VARCHAR(255) NOT NULL, 
    mime_type           VARCHAR(255) NOT NULL, 
    original_file_name  TEXT NOT NULL,
    description         TEXT,
    owner               UUID REFERENCES users(id) NOT NULL,
    status              VARCHAR(255) NOT NULL,
    refernce_count      INT DEFAULT 0,
    create_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    update_at           TIMESTAMPTZ
);

-- CREATE TABLE IF NOT EXISTS posts (
--     id                  UUID PRIMARY KEY,
--     author              UUID REFERENCES users(id) NOT NULL,
--     content             TEXT,
--     tag                 JSONB[],
--     attachement         JSONB[],
--     visibility          VARCHAR(255) NOT NULL,
--     -- group               UUID REFERENCES groups(id),
--     create_at           TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     update_at           TIMESTAMPTZ
-- );
-- 
-- CREATE TABLE IF NOT EXISTS posts (
--     id                  UUID PRIMARY KEY,
--     specific_user       UUID REFERENCES users(id) NOT NULL
-- );

