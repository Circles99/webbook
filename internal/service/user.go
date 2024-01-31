package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"webbook/internal/domain"
	"webbook/internal/respository"
)

var (
	ErrUserDuplicateEmail = respository.ErrUserDuplicateEmail
)

type UserService struct {
	repo *respository.UserRepository
}

func NewUserService(repo *respository.UserRepository) *UserService {
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
