package service

import (
	"context"
	"webbook/internal/domain"
	"webbook/internal/respository"
)

type UserService struct {
	repo *respository.UserRepository
}

func NewUserService(repo *respository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {

	return svc.repo.Create(ctx, u)
}
