package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webbook/internal/web"
	"webbook/ioc"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string
		// 提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name:   "发送成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)

			},
			reqBody:  `{"phone":15212345678}`,
			wantCode: 200,
			wantBody: web.Result{Code: 0, Msg: "发送成功"},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", time.Minute*9+time.Second*30).Result()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
			},
			reqBody:  `{"phone":15212345678}`,
			wantCode: 200,
			wantBody: web.Result{Code: 0, Msg: "发送太频繁，请稍后再试"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("content-type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)

			if resp.Code != 200 {
				return
			}

			var webResult web.Result
			err = json.NewDecoder(resp.Body).Decode(&webResult)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webResult)
			tc.after(t)
		})
	}
}
