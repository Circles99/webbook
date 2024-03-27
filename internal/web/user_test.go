package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webbook/internal/service"
	svcmocks "webbook/internal/service/mocks"
)

func TestUserHandler_Signup(t *testing.T) {

	var (
		tests = []struct {
			name     string
			mock     func(ctrl *gomock.Controller) service.UserService
			reqBody  string
			wantCode int
			wantBody string
		}{
			{
				name: "注册成功",
				mock: func(ctrl *gomock.Controller) service.UserService {
					usersvc := svcmocks.NewMockUserService(ctrl)
					usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
					return usersvc
				},
				reqBody: `
	{
		"email":"123@qq.com",
		"password":"123456"
	}
	`,
				wantBody: "注册成功",
				wantCode: 200,
			},
			{
				name: "参数不对， bind失败",
				mock: func(ctrl *gomock.Controller) service.UserService {
					usersvc := svcmocks.NewMockUserService(ctrl)
					usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
					return usersvc
				},
				wantCode: http.StatusBadRequest,
			},
			{
				name: "邮箱格式不对",
				mock: func(ctrl *gomock.Controller) service.UserService {
					usersvc := svcmocks.NewMockUserService(ctrl)
					return usersvc
				},
				reqBody: `
	{
		"email":"123@q",
		"password":"123456"
	}
	`,
				wantBody: "你的邮箱格式不对",
				wantCode: 200,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usersvc := svcmocks.NewMockUserService(ctrl)

			server := gin.Default()
			h := NewUserHandler(usersvc, nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tt.reqBody)))

			require.NoError(t, err)
			req.Header.Set("content-type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

}
