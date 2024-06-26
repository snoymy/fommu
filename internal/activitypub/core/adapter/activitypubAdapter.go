package adapter

import (
	"context"

	"github.com/go-ap/activitypub"
)

type ActivitypubAdapter interface {
    GetUserByUrl(ctx context.Context, url string) (*activitypub.Person, error)
}
