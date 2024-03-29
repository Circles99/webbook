package repository

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webbook/internal/domain"
	"webbook/internal/repository/cache"
	cachemocks "webbook/internal/repository/cache/mocks"
	"webbook/internal/repository/dao"
	daomocks "webbook/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.Unix())
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		ctx     context.Context
		id      int64
		want    domain.User
		wantErr error
	}{
		{
			name: "缓存未命中,查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrUserNotFound)

				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "15244448888",
						Valid:  true,
					},
					Password: "dsadsadsa",
					NickName: "xxxx",
					Created:  now.UnixMilli(),
					Updated:  now.UnixMilli(),
				}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "15244448888",
					Password: "dsadsadsa",
					NickName: "xxxx",
					Created:  now,
				}).Return(nil)

				return d, c
			},
			ctx: context.Background(),
			id:  123,
			want: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Phone:    "15244448888",
				Password: "dsadsadsa",
				NickName: "xxxx",
				Created:  now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中,查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "15244448888",
					Password: "dsadsadsa",
					NickName: "xxxx",
					Created:  now,
				}, nil)

				d := daomocks.NewMockUserDao(ctrl)

				return d, c
			},
			ctx: context.Background(),
			id:  123,
			want: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Phone:    "15244448888",
				Password: "dsadsadsa",
				NickName: "xxxx",
				Created:  now,
			},
			wantErr: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tt.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			user, err := repo.FindById(tt.ctx, tt.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, user)
			time.Sleep(time.Second)
		})
	}
}
