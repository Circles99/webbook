package respository

import (
	"context"
	"webbook/internal/domain"
)

type UserRepository struct {
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return nil
}
