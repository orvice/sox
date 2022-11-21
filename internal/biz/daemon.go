package biz

import (
	"context"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/mattn/go-mastodon"
	"github.com/orvice/sox/internal/config"
	"github.com/orvice/sox/pkg/storage"
	"github.com/weeon/contract"
	"github.com/weeon/log"
	"github.com/weeon/utils"
	"github.com/weeon/utils/ctxutil"
	"github.com/weeon/utils/task"
	"net/http"
	"time"
)

var(
	storages []storage.Storage
)

func Daemon(ctx context.Context) {
	task.NewTaskAndRun("tweet-sync", time.Second*15, func() error {
		checkTweets(ctx)
		return nil
	}, task.SetTaskLogger(log.GetLogger()))
}

func checkTweets(ctx context.Context) {
	ctx = ctxutil.AddRequestID(ctx, utils.NewUUID())
	ts, err := getTweets()
	if err != nil {
		log.Errorw("get tweets error",
			contract.ErrorMessage, err.Error(),
			contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
		)
		return
	}
	log.Infow("get tweets ",
		contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
		"len", len(ts),
		// "tweets",ts,
	)
	for _, t := range ts {
		checkTweet(ctx, t)
	}
}

func checkTweet(ctx context.Context, t twitter.Tweet) {
	var medias = make([]string, 0)
	var prefix string
	if t.Entities != nil && len(t.Entities.Media) != 0 {
		for _, m := range t.Entities.Media {
			medias = append(medias, m.MediaURLHttps)
		}
	}

	if t.InReplyToStatusID != 0 && config.Conf.SkipReply {
		return
	}

	if tm, err := t.CreatedAtTime(); err == nil {
		if tm.Before(time.Now().AddDate(0, 0, -1)) {
			log.Debugw("tweeted over one day, skip",
				"tweet_id", t.ID,
				contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
			)
			return
		}
	}

	if t.RetweetedStatus != nil {
		prefix = fmt.Sprintf("[RT @%s@twitter]", t.RetweetedStatus.User.ScreenName)
	}

	log.Debugf("prefix %s", prefix)

	for k, mcli := range mastodons {

		key := tweeID2Key(t.ID, k)
		has, err := cache.Has(ctx, key)
		if err != nil {
			log.Errorw("check tweet exist  fail",
				"tweet_id", t.ID,
				"key", key,
				contract.ErrorMessage, err.Error(),
				contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
				"media", medias,
			)
			continue
		}
		if has {
			log.Debugw("tweet is already sync to mastodon",
				"tweet_id", t.ID,
				"mastodon_id", k,
				contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
			)
			continue
		}

		mids := make([]mastodon.ID, 0)

		if len(medias) != 0 {
			for _, m := range medias {
				resp, err := http.DefaultClient.Get(m)
				if err != nil {
					log.Errorw("get media  fail",
						"tweet_id", t.ID,
						contract.ErrorMessage, err.Error(),
						contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
						"media", medias,
					)
					continue
				}

				att, err := mcli.UploadMediaFromReader(ctx, resp.Body)

				if err != nil {
					log.Errorw("upload media fail",
						"tweet_id", t.ID,
						contract.ErrorMessage, err.Error(),
						contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
						"media", medias,
					)
					continue
				}
				mids = append(mids, att.ID)
			}
		}
		toot := &mastodon.Toot{
			Status:      fmt.Sprintf("%s %s", prefix, t.Text),
			MediaIDs:    nil,
			Sensitive:   false,
			SpoilerText: "",
			Visibility:  "",
		}

		if len(mids) != 0 {
			toot.MediaIDs = mids
		}

		s, err := mcli.PostStatus(ctx, toot)
		if err != nil {
			log.Errorw("post toot fail",
				"tweet_id", t.ID,
				contract.ErrorMessage, err.Error(),
				contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
			)
			continue
		}

		err = cache.Write(ctx, key, []byte(s.ID))
		if err != nil {
			log.Errorw("write key error ",
				contract.ErrorMessage, err.Error(),
				contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
				"key", key,
			)
			continue
		}

		log.Infow("toot success ",
			contract.RequestID, ctxutil.GetRequestIDFromContext(ctx),
			"key", key,
			"toot", s,
		)

	}

}

func tweeID2Key(id int64, mid int) string {
	return fmt.Sprintf("%s-%d-%d", config.Conf.InstanceID, id, mid)
}
