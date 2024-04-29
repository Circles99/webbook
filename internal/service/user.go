package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webbook/internal/domain"
	"webbook/internal/repository"
	"webbook/pkg/logger"
)

var (
	ErrUserDuplicateEmail     = repository.ErrUserDuplicate
	ErrInvalidEmailOrPassword = errors.New("邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, userId int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	logger logger.Logger
}

func NewUserService(repo repository.UserRepository, l logger.Logger) UserService {
	return &UserServiceImpl{repo: repo, logger: l}
}

func (svc *UserServiceImpl) SignUp(ctx context.Context, u domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(password)
	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceImpl) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidEmailOrPassword
	}
	return u, nil
}

func (svc *UserServiceImpl) Edit(ctx context.Context, u domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return svc.repo.Edit(ctx, u)
}

func (svc *UserServiceImpl) Profile(ctx context.Context, userId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userId)
}

func (svc *UserServiceImpl) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断是否有这个用户

	if err != repository.ErrUserNotFound {
		return domain.User{}, err
	}

	svc.logger.Info("用户未注册", logger.String("phone", phone))
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})

	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里可能会遇到主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserServiceImpl) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenId)
	// 要判断是否有这个用户

	if err != repository.ErrUserNotFound {
		return u, err
	}

	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: info,
	})

	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里可能会遇到主从延迟问题
	return svc.repo.FindByWechat(ctx, info.OpenId)
}
