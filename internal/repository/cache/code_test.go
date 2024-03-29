package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webbook/internal/repository/cache/redismocks"
)

func TestRedisCodeCache_Set(t *testing.T) {

	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller, biz, phone, code string) redis.Cmdable
		ctx     context.Context
		biz     string
		code    string
		phone   string
		wantErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller, biz, phone, code string) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				//res.SetErr()
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return cmd
			},
			ctx:     nil,
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: nil,
		},

		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller, biz, phone, code string) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis"))
				//res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return cmd
			},
			ctx:     nil,
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: errors.New("mock redis"),
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller, biz, phone, code string) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				//res.SetErr(errors.New("mock redis"))
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return cmd
			},
			ctx:     nil,
			biz:     "login",
			code:    "123456",
			phone:   "152",
			wantErr: ErrSetCodeTooMany,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tt.mock(ctrl, tt.biz, tt.phone, tt.code))
			err := c.Set(tt.ctx, tt.biz, tt.phone, tt.code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
