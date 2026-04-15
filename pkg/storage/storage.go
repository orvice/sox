package storage

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mattn/go-mastodon"
)

type Storage interface {
	SaveTweet(ctx context.Context, t twitter.Tweet)
	SaveToot(ctx context.Context, t mastodon.Status)
}
