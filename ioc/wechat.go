package ioc

import "webbook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service {
	return wechat.NewWechatService("", "")
}
