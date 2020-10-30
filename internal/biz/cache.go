package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/orvice/twitter2mastodon/internal/config"
	"github.com/peterbourgon/diskv"
)

var (
	cache Cache
)

const (
	defaultExpire = time.Hour * 48
	driverRedis   = "redis"
)

type Cache interface {
	Has(ctx context.Context, key string) (bool, error)
	Write(ctx context.Context, key string, value []byte) error
}

func InitCache() error {
	var err error
	switch config.Conf.CacheDriver {
	case driverRedis:
		cache, err = NewRedisCache()
	default:
		cache, err = NewFileCache()
	}
	return err
}

type FileCache struct {
	d *diskv.Diskv
}

func NewFileCache() (*FileCache, error) {
	path := config.Conf.CachePath
	if path == "" {
		path = config.DefaultCachePath
	}
	flatTransform := func(s string) []string { return []string{} }
	// Initialize a new diskv store, rooted at "my-data-dir", with a 128MB cache.
	d := diskv.New(diskv.Options{
		BasePath:     path,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024 * 128,
	})
	err := d.Write("foo", []byte("bar"))
	if err != nil {
		return nil, err
	}
	return &FileCache{d: d}, nil
}

func (f *FileCache) Has(ctx context.Context, key string) (bool, error) {
	return f.d.Has(key), nil
}
func (f *FileCache) Write(ctx context.Context, key string, value []byte) error {
	return f.d.Write(key, value)
}

type RedisCache struct {
	cli *redis.Client
}

func NewRedisCache() (*RedisCache, error) {
	cfg := config.Conf.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})
	sc := rdb.Ping(context.Background())
	if sc.Err() != nil {
		return nil, sc.Err()
	}
	return &RedisCache{
		cli: rdb,
	}, nil
}

func (r *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	i, err := r.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return i > 0, nil
}
func (r *RedisCache) Write(ctx context.Context, key string, value []byte) error {
	return r.cli.Set(ctx, key, value, defaultExpire).Err()
}
