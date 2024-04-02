package ratelimit

import "context"

type Limiter interface {
	// Limited key是限流对象，是否触发限流
	Limit(ctx context.Context, key string) (bool, error)
}
