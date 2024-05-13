# Schema
```mermaid
erDiagram
USERS {
    id                      uuid            PK
    email                   varchar(255)    
    password_hash           varchar(255)
    status                  enum                    "posible value: active, inactive, deleted"
    username                varchar(255)
    display_name            varchar(255)
    name_prefix             varchar(255)
    name_suffix             varchar(255)
    bio                     text
    avartar                 varchar(255)            "url to image"
    banner                  varchar(255)            "url to image"
    tag                     array
    discoverable            boolean                 
    auto_approve_follower   boolean                 "default false"
    follower_count          int                     "defalult 0"
    following_count         int                     "defalult 0"
    public_key              text
    private_key             text
    url                     text
    remote                  boolean                 "default false"
    redirect_url            varchar(255)            "redirect_url, for user to redirect to new instance"
    create_at               timestamptz
    update_at               timestamptz
}
```

```mermaid
erDiagram
MEDIA {
    id                  uuid            PK
    url                 text
    type                varchar(255)
    description         text
    owner               uuid
    status              enum                        "posible value: visible, invisible, trash, deleted"
    refernce_count      int
    create_at           timestamptz
    update_at           timestamptz
}
```

# Relations
```mermaid
erDiagram
USERS ||--o{ SESSIONS: have
USERS ||--o{ POSTS: have
USERS ||--o{ POST_REACTS: have
USERS ||--o{ IMAGES: have
USERS ||--o{ ATTACHMENTS: have
USERS ||--o{ COMMENTS: have
USERS ||--o{ COMMENT_REACTS: have
USERS ||--o{ FOLLOWINGS: is
USERS ||--o{ MESSAGES: have
USERS ||--o{ CHATROOM_MEMBERS: is
USERS ||--o{ GROUP_MEMBERS: is

POSTS ||--o{ POST_REACTS: have
POSTS ||--o{ IMAGES: have
POSTS ||--o{ ATTACHMENTS: have
POSTS ||--o{ COMMENTS: have

COMMENTS ||--o{ COMMENTS: have
COMMENTS ||--o{ COMMENT_REACTS: have
COMMENTS ||--o{ IMAGES: have
COMMENTS ||--o{ ATTACHMENTS: have

GROUPS ||--o{ POSTS: have
GROUPS ||--o{ GROUP_MEMBERS: have
GROUPS ||--o{ GROUP_PERMISSIONS: have

GROUP_MEMBERS ||--o{ GROUP_PERMISSIONS: have

CHATROOMS ||--o{ CHATROOM_MEMBERS: have
CHATROOMS ||--o{ MESSAGES: have
```

# Refernce by

## USERS
```mermaid
erDiagram
USERS ||--o{ SESSIONS: have
USERS ||--o{ POSTS: have
USERS ||--o{ POST_REACTS: have
USERS ||--o{ IMAGES: have
USERS ||--o{ ATTACHMENTS: have
USERS ||--o{ COMMENTS: have
USERS ||--o{ COMMENT_REACTS: have
USERS ||--o{ FOLLOWINGS: is
USERS ||--o{ MESSAGES: have
USERS ||--o{ CHATROOM_MEMBERS: is
USERS ||--o{ GROUP_MEMBERS: is
```

## POSTS
```mermaid
erDiagram
POSTS ||--o{ POST_REACTS: have
POSTS ||--o{ IMAGES: have
POSTS ||--o{ ATTACHMENTS: have
POSTS ||--o{ COMMENTS: have
```

## COMMENTS
```mermaid
erDiagram
COMMENTS ||--o{ COMMENTS: have
COMMENTS ||--o{ COMMENT_REACTS: have
COMMENTS ||--o{ IMAGES: have
COMMENTS ||--o{ ATTACHMENTS: have
```

## GROUPS
```mermaid
erDiagram
GROUPS ||--o{ POSTS: have
GROUPS ||--o{ GROUP_MEMBERS: have
GROUPS ||--o{ GROUP_PERMISSIONS: have
```

## GROUP_MEMBERS
```mermaid
erDiagram
GROUP_MEMBERS ||--o{ GROUP_PERMISSIONS: have
```

## CHATROOMS
```mermaid
erDiagram
CHATROOMS ||--o{ CHATROOM_MEMBERS: have
CHATROOMS ||--o{ MESSAGES: have
```

# Refernce

## SESSIONS
```mermaid
erDiagram
USERS ||--o{ SESSIONS: have
```

## POSTS
```mermaid
erDiagram
USERS ||--o{ POSTS: have
GROUPS ||--o{ POSTS: have
```

## POST_REACTS
```mermaid
erDiagram
USERS ||--o{ POST_REACTS: have
POSTS ||--o{ POST_REACTS: have
```

## IMAGES
```mermaid
erDiagram
USERS ||--o{ IMAGES: have
POSTS ||--o{ IMAGES: have
COMMENTS ||--o{ IMAGES: have
```

## ATTACHMENTS
```mermaid
erDiagram
USERS ||--o{ ATTACHMENTS: have
POSTS ||--o{ ATTACHMENTS: have
COMMENTS ||--o{ ATTACHMENTS: have
```

## COMMENTS
```mermaid
erDiagram
USERS ||--o{ COMMENTS: have
POSTS ||--o{ COMMENTS: have
COMMENTS ||--o{ COMMENTS: have
```

## COMMENT_REACTS
```mermaid
erDiagram
USERS ||--o{ COMMENT_REACTS: have
COMMENTS ||--o{ COMMENT_REACTS: have
```

## FOLLOWINGS
```mermaid
erDiagram
USERS ||--o{ FOLLOWINGS: is
```

## MESSAGES
```mermaid
erDiagram
USERS ||--o{ MESSAGES: have
CHATROOMS ||--o{ MESSAGES: have
```

## CHATROOM_MEMBERS
```mermaid
erDiagram
USERS ||--o{ CHATROOM_MEMBERS: is
CHATROOMS ||--o{ CHATROOM_MEMBERS: have
```

## GROUP_MEMBERS
```mermaid
erDiagram
USERS ||--o{ GROUP_MEMBERS: is
GROUPS ||--o{ GROUP_MEMBERS: have
```

## GROUP_PERMISSIONS
```mermaid
erDiagram
GROUPS ||--o{ GROUP_PERMISSIONS: have
GROUP_MEMBERS ||--o{ GROUP_PERMISSIONS: have
```
