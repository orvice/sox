package biz

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/orvice/sox/internal/config"
)

var (
	twitterCli *twitter.Client
)

func InitTwitter() error {
	cfg := config.Conf
	conf := oauth1.NewConfig(cfg.TwitterKey, cfg.TwitterSecretKey)
	token := oauth1.NewToken(cfg.TwitterToken, cfg.TwitterTokenSecret)
	httpClient := conf.Client(oauth1.NoContext, token)

	// Twitter client
	twitterCli = twitter.NewClient(httpClient)
	return nil
}

func getTweets() ([]twitter.Tweet, error) {
	ts, resp, err := twitterCli.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: config.Conf.TwitterScreenName,
	})

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return ts, err
	}

	return ts, nil
}
