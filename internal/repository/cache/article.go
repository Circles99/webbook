package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webbook/internal/domain"
)

// 高并发情况下提前预加载+本地缓存+一致性哈希负载均衡算法

type ArticleCache interface {
	GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, userId int64) error
	Set(ctx context.Context, id int64, art domain.Article) error
	SetPub(ctx context.Context, art domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) SetPub(ctx context.Context, art domain.Article) error {
	return nil
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error) {
	bs, err := r.client.Get(ctx, r.FirstKey(userId)).Bytes()
	if err != nil {
		return nil, err
	}

	var arts []domain.Article
	err = json.Unmarshal(bs, &arts)

	return arts, err
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.FirstKey(userId), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, userId int64) error {
	return r.client.Del(ctx, r.FirstKey(userId)).Err()
}

func (r *RedisArticleCache) Set(ctx context.Context, id int64, art domain.Article) error {

	bs, err := json.Marshal(art)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.key(id), bs, time.Second*10).Err()
}

func (r *RedisArticleCache) key(userId int64) string {
	return fmt.Sprintf("article:%d", userId)
}

func (r *RedisArticleCache) FirstKey(userId int64) string {
	return fmt.Sprintf("article_first_page:%d", userId)
}
