package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webbook/internal/domain"
	"webbook/internal/service"
	svcmocks "webbook/internal/service/mocks"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/ginx"
	"webbook/pkg/logger"
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
				wantBody: "邮箱错误",
				wantCode: 200,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tt.mock(ctrl), nil, nil, nil)
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

func TestUserHandler_LogSms(t *testing.T) {

	var (
		tests = []struct {
			name     string
			mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService, logger.Logger, ijwt.Handler)
			reqBody  string
			wantCode int
			wantBody ginx.Result
		}{
			{
				name: "登录成功",
				mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, logger.Logger, ijwt.Handler) {
					codeSvc := svcmocks.NewMockCodeService(ctrl)
					codeSvc.EXPECT().Verify(gomock.Any(), "login", "15211112222", "123456").Return(true, nil)

					usersvc := svcmocks.NewMockUserService(ctrl)
					usersvc.EXPECT().FindOrCreate(gomock.Any(), "15211112222").Return(domain.User{
						Id:    123,
						Email: "123@qq.com",
						Phone: "15211112222",
					}, nil)
					return usersvc, codeSvc, nil, nil
				},
				reqBody:  `{"phone":"15211112222", "code":"123456"}`,
				wantCode: http.StatusOK,
				wantBody: ginx.Result{
					Code: 4,
					Msg:  "验证成功",
				},
			},
			{
				name: "登录失败，验证码错误",
				mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, logger.Logger, ijwt.Handler) {
					codeSvc := svcmocks.NewMockCodeService(ctrl)
					codeSvc.EXPECT().Verify(gomock.Any(), "login", "15211112222", "123456").Return(false, nil)

					usersvc := svcmocks.NewMockUserService(ctrl)
					return usersvc, codeSvc, nil, nil
				},
				reqBody:  `{"phone":"15211112222", "code":"123456"}`,
				wantCode: http.StatusOK,
				wantBody: ginx.Result{
					Code: 4,
					Msg:  "验证码有误",
				},
			},

			{
				name: "登录失败，系统错误",
				mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService, logger.Logger, ijwt.Handler) {
					codeSvc := svcmocks.NewMockCodeService(ctrl)
					codeSvc.EXPECT().Verify(gomock.Any(), "login", "15211112222", "123456").Return(false, errors.New("mock 错误"))

					usersvc := svcmocks.NewMockUserService(ctrl)

					return usersvc, codeSvc, nil, nil
				},
				reqBody:  `{"phone":"15211112222", "code":"123456"}`,
				wantCode: http.StatusOK,
				wantBody: ginx.Result{
					Code: 5,
					Msg:  "系统错误",
				},
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tt.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tt.reqBody)))

			require.NoError(t, err)
			req.Header.Set("content-type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)

			var res ginx.Result
			err = json.NewDecoder(resp.Body).Decode(&res)

			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, res)

		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)

	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

}
