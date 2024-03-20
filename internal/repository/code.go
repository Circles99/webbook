package repository

import (
	"context"
	"webbook/internal/repository/cache"
)

var (
	ErrSetCodeTooMany = cache.ErrSetCodeTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{cache: c}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
