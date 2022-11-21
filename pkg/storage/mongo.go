package storage

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mattn/go-mastodon"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStr struct {
	client *mongo.Client
}

func (m *MongoStr) SaveTweet(ctx context.Context, t twitter.Tweet) {

}

func (m *MongoStr) SaveToot(ctx context.Context, t mastodon.Status) {

}
