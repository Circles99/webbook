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
	IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	// AddCollectionItem 收藏
	AddCollectionItem(ctx context.Context, biz string, bizId, cid, uid int64) error
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

func (c *CachedReadCntRepository) IncrLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	// 先插入点赞，然后更新点赞
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}

	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)

}

func (c *CachedReadCntRepository) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}

	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid, uid int64) error {

	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Cid:   cid,
		BizId: bizId,
		Biz:   biz,
		Uid:   uid,
	})
	if err != nil {
		return err
	}

	// 收藏个数

	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) Get(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
	//TODO implement me
	panic("implement me")
}
