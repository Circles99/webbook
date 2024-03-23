package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webbook/internal/domain"
	"webbook/internal/repository"
)

var (
	ErrUserDuplicateEmail     = repository.ErrUserDuplicate
	ErrInvalidEmailOrPassword = errors.New("邮箱或密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(password)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
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

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return svc.repo.Edit(ctx, u)
}

func (svc *UserService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userId)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断是否有这个用户

	if err != repository.ErrUserNotFound {
		return domain.User{}, err
	}

	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})

	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里可能会遇到主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}
