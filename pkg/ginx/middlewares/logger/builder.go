package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	logFunc       func(ctx context.Context, al *AccessLog)
	allowReqBody  bool
	allowRespBody bool
}

type AccessLog struct {
	Method     string
	Url        string
	Duration   string
	ReqBody    string
	RespBody   string
	StatusCode int
}

func NewMiddlewareBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: fn,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			Url:    url,
		}

		if b.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			// body 是一个ReadCloser(流)，读完就没了，需要重新返回去
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			// 会引起复制，很消耗CPU和内存
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			c.Writer = responseWriter{
				al:             al,
				ResponseWriter: c.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()

			b.logFunc(c, al)
		}()

		// 执行业务逻辑
		c.Next()

	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}
