package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"webbook/internal/service/sms"
)

type TokenAuthSmsService struct {
	svc sms.Service
	key string
}

// Send 发送， 其中biz为线下申请的业务token
func (t TokenAuthSmsService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

	var claims TokenClaims
	// 如果这里能解析成功，说明就是对应的业务方
	token, err := jwt.ParseWithClaims(biz, &claims, func(token *jwt.Token) (interface{}, error) {
		return t.key, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token 不合法")
	}

	return t.svc.Send(ctx, claims.Tpl, args, numbers...)
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Tpl string
}
