package ioc

import (
	"webbook/internal/service/oauth2/wechat"
	"webbook/internal/web"
)

func InitOAuth2WechatService() wechat.Service {
	return wechat.NewWechatService("", "")
}

func NewWechatHandlerConfig() web.WechatConfig {
	return web.WechatConfig{
		Secure: true,
	}
}
