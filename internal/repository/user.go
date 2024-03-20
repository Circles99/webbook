package repository

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{dao: dao, cache: c}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这操作缓存以及其他操作
}

func (r *UserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Desc:     u.Desc,
	})
	// 在这操作缓存以及其他操作
}

func (r *UserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {

	u, err := r.cache.Get(ctx, userId)
	if err == nil {
		return u, nil
	}

	//if err == cache.ErrUserNotFound {
	//
	//}
	user, err := r.dao.FindById(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}

	u = domain.User{
		Id:       user.Id,
		Email:    user.Email,
		NickName: user.NickName,
		Birthday: user.Birthday,
		Desc:     user.Desc,
	}
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 写入日志
		}
	}()

	return u, nil
}
