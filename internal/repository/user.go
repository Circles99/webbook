package repository

import (
	"context"
	"database/sql"
	"time"
	"webbook/internal/domain"
	"webbook/internal/repository/cache"
	"webbook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	Edit(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, userId int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{dao: dao, cache: c}
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
	// 在这操作缓存以及其他操作
}

func (r *CacheUserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, r.domainToEntity(u))
	// 在这操作缓存以及其他操作
}

func (r *CacheUserRepository) FindById(ctx context.Context, userId int64) (domain.User, error) {

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

	u = r.entityToDomain(user)
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 写入日志
		}
	}()

	return u, nil
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)

	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		NickName: u.NickName,
		Birthday: u.Birthday,
		Desc:     u.Desc,
		Created:  time.UnixMilli(u.Created),
	}
}

func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		NickName: u.NickName,
		Birthday: u.Birthday,
		Desc:     u.Desc,
	}
}
