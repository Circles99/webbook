package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"time"
	"webbook/internal/domain"
	"webbook/internal/repository"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao/article"
	"webbook/pkg/logger"
)

type ArticleRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, authorId int64, status domain.ArticleStatus) error
	List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)
	PubDetail(ctx context.Context, id int64) (domain.Article, error)
}

type ArticleRepositoryImpl struct {
	userRepo repository.UserRepository
	dao      article.ArticleDAO
	cache    cache.ArticleCache
	l        logger.Logger
	//// v1 操作两个dao
	//authorDao article.AuthorDao
	//readerDao article.ReaderDao
}

func NewArticleRepository(dao article.ArticleDAO, cache cache.ArticleCache, l logger.Logger) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (a *ArticleRepositoryImpl) PubDetail(ctx context.Context, id int64) (domain.Article, error) {

	// 读取线上库数据
	art, err := a.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	// 组装user
	user, err := a.userRepo.FindById(ctx, art.AuthorId)
	res := a.toDomain(article.Article(art))
	res.Author.Name = user.NickName
	return res, nil

}

func (a *ArticleRepositoryImpl) Detail(ctx context.Context, id int64) (domain.Article, error) {
	data, err := a.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return a.toDomain(data), nil
}

func (a *ArticleRepositoryImpl) List(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error) {

	if offset == 0 && limit <= 100 {
		data, err := a.cache.GetFirstPage(ctx, userId)
		if err == nil {
			go func() {
				// 预加载缓存第一篇文章
				a.preCache(ctx, data)
			}()
			return data[:limit], nil
		}

	}

	res, err := a.dao.GetByAuthor(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}

	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return a.toDomain(src)
	})
	go func() {
		if offset == 0 && limit <= 100 {
			err = a.cache.SetFirstPage(ctx, userId, data)
			if err != nil {
				a.l.Error("回写缓存失败", logger.Error(err))
			}
			// 预加载缓存第一篇文章
			a.preCache(ctx, data)
		}
	}()

	return data, nil
}

func (a *ArticleRepositoryImpl) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author:  domain.Author{Id: art.AuthorId},
		Status:  domain.ArticleStatus(art.Status),
		Created: time.UnixMilli(art.Created),
		Updated: time.UnixMilli(art.Updated),
	}
}

func (a *ArticleRepositoryImpl) SyncStatus(ctx context.Context, id, authorId int64, status domain.ArticleStatus) error {
	defer func() {
		// 删除缓存
		err := a.cache.DelFirstPage(ctx, authorId)
		if err != nil {
			a.l.Error("删除缓存失败", logger.Error(err))
		}
	}()
	return a.dao.SyncStatus(ctx, id, authorId, status.ToUint8())
}

// dao 层面上解决事务问题
func (a *ArticleRepositoryImpl) Sync(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		// 删除缓存
		err := a.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			a.l.Error("删除缓存失败", logger.Error(err))
		}

	}()
	id, err := a.dao.Sync(ctx, a.toEntity(art))
	if err == nil {
		// 删除缓存
		er := a.cache.DelFirstPage(ctx, art.Author.Id)
		if er != nil {
			a.l.Error("删除缓存失败", logger.Error(er))
		}
		// 提前缓存线上库数据
		er = a.cache.SetPub(ctx, art)
		if er != nil {
			a.l.Warn("设置线上库失败", logger.Error(er))
		}
	}
	return id, err
}

//// Repository 层面上解决事务问题
//func (a *ArticleRepositoryImpl) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
//	// 先保存到制作库，在保存到线上库
//
//	var (
//		id  = art.Id
//		err error
//	)
//
//	if art.Id > 0 {
//		err = a.authorDao.Update(ctx, a.toEntity(art))
//	} else {
//		id, err = a.authorDao.Insert(ctx, a.toEntity(art))
//	}
//	if err != nil {
//		return id, err
//	}
//	// 操作线上库，保存数据
//	// 考虑到，线上可能有 可能没有，要有一个upset的写法
//	id, err = a.readerDao.Upsert(ctx, a.toEntity(art))
//
//	return id, err
//}

func (a *ArticleRepositoryImpl) Save(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		// 删除缓存
		err := a.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			a.l.Error("删除缓存失败", logger.Error(err))
		}
	}()

	return a.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (a *ArticleRepositoryImpl) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		// 删除缓存
		err := a.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			a.l.Error("删除缓存失败", logger.Error(err))
		}
	}()
	return a.dao.Update(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	})
}

func (a *ArticleRepositoryImpl) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func (a *ArticleRepositoryImpl) preCache(ctx context.Context, data []domain.Article) {
	// 不预加载缓存大对象
	if len(data) > 0 && len(data[0].Content) <= 1024*1024 {
		err := a.cache.Set(ctx, data[0].Id, data[0])
		if err != nil {
			a.l.Error("提前预加载缓存失败", logger.Error(err))
		}
	}
}
