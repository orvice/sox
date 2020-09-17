package biz

import (
	"github.com/mattn/go-mastodon"
	"github.com/orvice/twitter2mastodon/internal/config"
)

var (
	mastodons []*mastodon.Client
)

func InitMastodon() error {
	for _, v := range config.Conf.Mastodons {
		if !v.Enable {
			continue
		}
		c := mastodon.NewClient(&mastodon.Config{
			Server:       v.Server,
			ClientID:     v.ClientID,
			ClientSecret: v.ClientSecret,
			AccessToken:  v.Token,
		})
		mastodons = append(mastodons, c)
	}
	return nil
}
