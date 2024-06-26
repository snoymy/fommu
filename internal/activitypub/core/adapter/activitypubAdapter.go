package adapter

import (
	"context"

	"github.com/snoymy/activitypub"
)

type ActivitypubAdapter interface {
    GetUserByUrl(ctx context.Context, url string) (*activitypub.Person, error)
}
