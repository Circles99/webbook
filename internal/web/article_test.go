package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webbook/internal/service"
	ijwt "webbook/internal/web/jwt"
	"webbook/pkg/logger"
)

func TestArticleHandler_Publish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService, logger.Logger, ijwt.Handler)
		reqBody  string
		wantCode int
		wantBody Result
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tt.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/article/publish", bytes.NewBuffer([]byte(tt.reqBody)))

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
