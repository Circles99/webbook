package repository

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
	"webbook/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// Collect 收藏
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDao
	l     logger.Logger
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}

	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedReadCntRepository) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedReadCntRepository) Collect(ctx context.Context, biz string, bizId, cid, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedReadCntRepository) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}
