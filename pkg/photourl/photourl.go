package photourl

import (
	"context"
	"time"
)

type Resolver interface {
	PublicURL(key string) string
	PresignedGetURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type Options struct {
	UsePresigned bool
	PresignTTL   time.Duration
}

// Resolve - возвращает готовый URL или nil(если ключ пустой/ошибка presign)
func Resolve(ctx context.Context, r Resolver, photoKey *string, opt Options) *string {
	if r == nil || photoKey == nil || *photoKey == "" {
		return nil
	}

	if opt.UsePresigned {
		u, err := r.PresignedGetURL(ctx, *photoKey, opt.PresignTTL)
		if err != nil {
			return nil
		}
		return &u
	}

	u := r.PublicURL(*photoKey)
	return &u
}
