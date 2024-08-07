package httpclient

import (
	"context"

	"github.com/snoymy/activitypub"
)

type ActivitypubClient interface {
	// FetchWebfinger retrieves the Webfinger response for the given URLs.
	FetchWebfinger(ctx context.Context, domain string, username string) ([]interface{}, error)
	
	// // FetchActivitypubObject retrieves a general ActivityPub object from the specified URL.
    FetchActor(ctx context.Context, url string) (*activitypub.Actor, error)
	// 
	// // FetchActivitypubLink retrieves an ActivityPub Link object from the specified URL.
	// FetchActivitypubLink(ctx context.Context, url string) (*activitypub.Link, error)
	// 
	// // FetchActivitypubActivity retrieves an ActivityPub Activity object from the specified URL.
	// FetchActivitypubActivity(ctx context.Context, url string) (*activitypub.Activity, error)
	// 
	// // FetchActivitypubIntransitiveActivity retrieves an ActivityPub Intransitive Activity object from the specified URL.
	// FetchActivitypubIntransitiveActivity(ctx context.Context, url string) (*activitypub.IntransitiveActivity, error)
	// 
	// // FetchActivitypubCollection retrieves an ActivityPub Collection object from the specified URL.
	// FetchActivitypubCollection(ctx context.Context, url string) (*activitypub.Collection, error)
	// 
	// // FetchActivitypubOrderedCollection retrieves an ActivityPub Ordered Collection object from the specified URL.
	// FetchActivitypubOrderedCollection(ctx context.Context, url string) (*activitypub.OrderedCollection, error)
	// 
	// // FetchActivitypubCollectionPage retrieves an ActivityPub Collection Page object from the specified URL.
	// FetchActivitypubCollectionPage(ctx context.Context, url string) (*activitypub.CollectionPage, error)
	// 
	// FetchActivitypubOrderedCollectionPage retrieves an ActivityPub Ordered Collection Page object from the specified URL.
	FetchOrderedCollectionPage(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error)
}

