package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webbook/internal/domain"
	"webbook/internal/repository"
	repmocks "webbook/internal/repository/mocks"
	"webbook/pkg/logger"
)

func TestUserServiceImpl_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$.re6u5Sud/V/fMuPYnLNje4Ha3BoMW58o.UwojIfFAzAEdJUhYwVq",
					Phone:    "152618925xx",
					Created:  now,
				}, nil)

				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$.re6u5Sud/V/fMuPYnLNje4Ha3BoMW58o.UwojIfFAzAEdJUhYwVq",
				Phone:    "152618925xx",
				Created:  now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)

				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{},
			wantErr:  ErrInvalidEmailOrPassword,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("mock 错误"))

				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{},
			wantErr:  errors.New("mock 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repmocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$.re6u5Sud/V/fMuPYnLNje4Ha3BoMW58o.UwojIfFAzAEdJUhYwVq",
					Phone:    "152618925xx",
					Created:  now,
				}, nil)

				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "1234562",
			wantUser: domain.User{},
			wantErr:  ErrInvalidEmailOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl), &logger.ZapLogger{})

			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
