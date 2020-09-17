package biz

import (
	"github.com/orvice/twitter2mastodon/internal/config"
	"github.com/peterbourgon/diskv"
)

var (
	dk Cache
)

type Cache interface {
	Has(key string) bool
	Write(key string, value []byte) error
}

func InitCache() error {
	path := config.Conf.CachePath
	if path == "" {
		path = config.DefaultCachePath
	}
	flatTransform := func(s string) []string { return []string{} }
	// Initialize a new diskv store, rooted at "my-data-dir", with a 128MB cache.
	dk = diskv.New(diskv.Options{
		BasePath:     path,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024 * 128,
	})

	err := dk.Write("foo", []byte("bar"))
	if err != nil {
		return err
	}

	return nil
}
