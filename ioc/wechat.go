package ioc

import (
	"webbook/internal/service/oauth2/wechat"
	"webbook/internal/web"
	"webbook/pkg/logger"
)

func InitOAuth2WechatService(l logger.Logger) wechat.Service {
	return wechat.NewWechatService("", "", l)
}

func NewWechatHandlerConfig() web.WechatConfig {
	return web.WechatConfig{
		Secure: true,
	}
}
